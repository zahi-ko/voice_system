package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"voice_system/internal/adapters/audioio"
	"voice_system/internal/adapters/memory"
	"voice_system/internal/application"

	"github.com/labstack/echo/v5"
)

// newEffectsRouter builds a fully wired router backed by a fresh in-memory
// store, with test.wav already uploaded. Returns the router and the new ID.
func newEffectsRouter(t *testing.T) (*echo.Echo, string) {
	t.Helper()

	fixture := filepath.Join("..", "..", "..", "..", "tests", "data", "test.wav")
	content, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}

	store, err := memory.NewAudioStore()
	if err != nil {
		t.Fatal(err)
	}
	registry, err := memory.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(application.NewAudioService(audioio.NewDecoder(), audioio.NewEncoder(), store), application.NewEffectService(store, registry), application.NewAnalysisService(store))
	router := echo.New()
	Register(router, handler)

	uploadBody, contentType := multipartFixture(t, "test.wav", content)
	uploadRequest := httptest.NewRequest(http.MethodPost, "/audio/upload", uploadBody)
	uploadRequest.Header.Set(echo.HeaderContentType, contentType)
	uploadResponse := httptest.NewRecorder()
	router.ServeHTTP(uploadResponse, uploadRequest)
	if uploadResponse.Code != http.StatusCreated {
		t.Fatalf("upload status = %d, want %d: %s", uploadResponse.Code, http.StatusCreated, uploadResponse.Body.String())
	}

	var uploaded struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(uploadResponse.Body.Bytes(), &uploaded); err != nil {
		t.Fatalf("decode upload response: %v", err)
	}
	return router, uploaded.ID
}

func TestListEffects(t *testing.T) {
	router, _ := newEffectsRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/effects/list", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var names []string
	if err := json.Unmarshal(rec.Body.Bytes(), &names); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	want := []string{"gain", "tempo", "reverse", "tempo_pv", "normalize"}
	for _, w := range want {
		found := false
		for _, n := range names {
			if n == w {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("effects list %v missing %q", names, w)
		}
	}
}

func TestApplyEffectReverse(t *testing.T) {
	router, id := newEffectsRouter(t)

	body := []byte(`{"name":"reverse"}`)
	req := httptest.NewRequest(http.MethodPost, "/effects/"+id+"/apply", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("apply status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var applied struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Duration float32 `json:"duration"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &applied); err != nil {
		t.Fatalf("decode apply response: %v", err)
	}
	if applied.ID != id {
		t.Errorf("apply response id = %q, want %q (saveAsNew 默认 false，应原地覆盖)", applied.ID, id)
	}
	if !strings.HasSuffix(applied.Name, "_reverse") {
		t.Errorf("apply response name = %q, want _reverse suffix", applied.Name)
	}
}

// TestApplyEffectTempo 覆盖参数必须嵌套在 parameters 下的契约：
// 平铺参数会被静默丢弃，导致 tempo 取到零值并被范围校验拒绝。
func TestApplyEffectTempo(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		wantRate   float64
		wantSuffix string
	}{
		{"tempo", `{"name":"tempo","parameters":{"tempo":1.5}}`, 66150, "_tempo"},
		{"tempo_pv", `{"name":"tempo_pv","parameters":{"tempo":1.5}}`, 44100, "_tempo_pv"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router, id := newEffectsRouter(t)

			req := httptest.NewRequest(http.MethodPost, "/effects/"+id+"/apply", strings.NewReader(c.body))
			req.Header.Set(echo.HeaderContentType, "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("apply status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
			}

			var applied struct {
				ID         string  `json:"id"`
				Name       string  `json:"name"`
				Duration   float32 `json:"duration"`
				SampleRate float64 `json:"sample_rate"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &applied); err != nil {
				t.Fatalf("decode apply response: %v", err)
			}
			if !strings.HasSuffix(applied.Name, c.wantSuffix) {
				t.Errorf("apply response name = %q, want %s suffix", applied.Name, c.wantSuffix)
			}
			if applied.SampleRate != c.wantRate {
				t.Errorf("sample_rate = %v, want %v", applied.SampleRate, c.wantRate)
			}
			// 1.5 倍速后时长应缩短为原来的 2/3
			if applied.Duration > 12.78 || applied.Duration < 8.0 {
				t.Errorf("duration = %f, want ~8.52 (原 12.78 / 1.5)", applied.Duration)
			}
		})
	}
}

func TestApplyEffectInvalidPayload(t *testing.T) {
	router, id := newEffectsRouter(t)

	cases := []struct {
		name string
		body string
	}{
		{"missing name", `{}`},
		{"unknown effect", `{"name":"pitch"}`},
		{"gain out of range", `{"name":"gain","parameters":{"db":999}}`},
		{"tempo out of range", `{"name":"tempo","parameters":{"tempo":4.0}}`},
		// 参数平铺在顶层等价于缺失 parameters，tempo 会取到零值并被拒绝
		{"tempo flat params", `{"name":"tempo","tempo":1.5}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/effects/"+id+"/apply", strings.NewReader(c.body))
			req.Header.Set(echo.HeaderContentType, "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// 效果名/参数问题属于调用方错误，应回 400 并带出具体原因
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("apply status = %d, want %d: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), "effects:") {
				t.Errorf("error body = %s, want it to carry the underlying reason", rec.Body.String())
			}
		})
	}
}

func TestApplyEffectChain(t *testing.T) {
	router, id := newEffectsRouter(t)

	body := []byte(`[{"name":"reverse"},{"name":"reverse"}]`)
	req := httptest.NewRequest(http.MethodPost, "/effects/"+id+"/chain", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("chain status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var response struct {
		ID             string `json:"id"`
		AppliedEffects string `json:"applied_effects"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode chain response: %v", err)
	}
	if response.ID != id {
		t.Errorf("chain response id = %q, want %q (saveAsNew 默认 false，应原地覆盖)", response.ID, id)
	}
	if response.AppliedEffects != "reverse->reverse->" {
		t.Errorf("applied_effects = %q, want %q", response.AppliedEffects, "reverse->reverse->")
	}
}

func TestApplyEffectChainEmpty(t *testing.T) {
	router, id := newEffectsRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/effects/"+id+"/chain", strings.NewReader(`[]`))
	req.Header.Set(echo.HeaderContentType, "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("empty chain status = %d, want %d: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestApplyEffectChainNameTruncated(t *testing.T) {
	router, id := newEffectsRouter(t)

	// 40 次 reverse 的链名远超 64 字节上限，应被截断为头尾保留形式
	chain := make([]string, 40)
	for i := range chain {
		chain[i] = `{"name":"reverse"}`
	}
	body := []byte("[" + strings.Join(chain, ",") + "]")
	req := httptest.NewRequest(http.MethodPost, "/effects/"+id+"/chain", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("chain status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var response struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode chain response: %v", err)
	}
	if len(response.Metadata.Name) > 64 {
		t.Errorf("chain name length = %d, want <= 64: %q", len(response.Metadata.Name), response.Metadata.Name)
	}
	if !strings.Contains(response.Metadata.Name, "...") {
		t.Errorf("truncated chain name %q should contain \"...\"", response.Metadata.Name)
	}
}
