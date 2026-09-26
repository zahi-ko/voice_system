package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"voice_system/internal/adapters/audioio"
	"voice_system/internal/adapters/memory"
	"voice_system/internal/application"
	"voice_system/internal/domain/audio"
	transporthttp "voice_system/internal/transport/http"
	"voice_system/internal/transport/http/generated"

	"github.com/labstack/echo/v5"
)

func TestAudioEndpointsSupportFixtureFormats(t *testing.T) {
	fixtures, err := filepath.Glob(filepath.Join("..", "..", "..", "..", "tests", "data", "test.*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(fixtures) == 0 {
		t.Fatal("no audio fixtures found")
	}

	for _, fixture := range fixtures {
		t.Run(filepath.Ext(fixture), func(t *testing.T) {
			content, err := os.ReadFile(fixture)
			if err != nil {
				t.Fatal(err)
			}

			store, err := memory.NewAudioStore()
			if err != nil {
				t.Fatal(err)
			}
			handler := NewHandler(application.NewAudioService(audioio.NewDecoder(), audioio.NewEncoder(), store))
			router := echo.New()
			Register(router, handler)

			uploadBody, contentType := multipartFixture(t, filepath.Base(fixture), content)
			uploadRequest := httptest.NewRequest(http.MethodPost, "/audio/upload", uploadBody)
			uploadRequest.Header.Set(echo.HeaderContentType, contentType)
			uploadResponse := httptest.NewRecorder()
			router.ServeHTTP(uploadResponse, uploadRequest)

			if uploadResponse.Code != http.StatusCreated {
				t.Fatalf("upload status = %d, want %d: %s", uploadResponse.Code, http.StatusCreated, uploadResponse.Body.String())
			}

			var uploaded generated.AudioMeta
			if err := json.Unmarshal(uploadResponse.Body.Bytes(), &uploaded); err != nil {
				t.Fatalf("decode upload response: %v", err)
			}
			if uploaded.ID == nil || *uploaded.ID == "" {
				t.Fatal("upload response has no audio ID")
			}
			wantName := strings.TrimSuffix(filepath.Base(fixture), filepath.Ext(fixture))
			if uploaded.Name == nil || *uploaded.Name != wantName {
				t.Fatalf("upload response name = %v, want %s", uploaded.Name, wantName)
			}

			listRequest := httptest.NewRequest(http.MethodGet, "/audio", nil)
			listResponse := httptest.NewRecorder()
			router.ServeHTTP(listResponse, listRequest)
			if listResponse.Code != http.StatusOK {
				t.Fatalf("list status = %d, want %d", listResponse.Code, http.StatusOK)
			}
			var listed generated.AudioList
			if err := json.Unmarshal(listResponse.Body.Bytes(), &listed); err != nil {
				t.Fatalf("decode list response: %v", err)
			}
			if len(listed.Items) != 1 || listed.Items[0].ID == nil || *listed.Items[0].ID != *uploaded.ID {
				t.Fatalf("list response = %+v, want uploaded audio", listed)
			}

			downloadRequest := httptest.NewRequest(http.MethodGet, "/audio/"+*uploaded.ID+"/download", nil)
			downloadResponse := httptest.NewRecorder()
			router.ServeHTTP(downloadResponse, downloadRequest)
			if downloadResponse.Code != http.StatusOK {
				t.Fatalf("download status = %d, want %d", downloadResponse.Code, http.StatusOK)
			}
			// 上传时音频已解码为归一化 float32 单声道，下载响应是重新编码的产物，
			// 字节流不可能与原始 fixture 相同（mp3 有损，且重编码为双声道）。
			// 语义上的往返无损 = 重新编码后的音频能被解回同样的采样率与时长
			// （解码端会下混回 mono，样本数按采样率换算即为时长）。
			expectedMIME := transporthttp.FormatToMIME(audio.Format(strings.TrimPrefix(filepath.Ext(fixture), ".")))
			if got := downloadResponse.Header().Get(echo.HeaderContentType); got != expectedMIME {
				t.Fatalf("download content-type = %q, want %q", got, expectedMIME)
			}

			downloadedMeta, err := audioio.NewDecoder().DecodeMeta(bytes.NewReader(downloadResponse.Body.Bytes()))
			if err != nil {
				t.Fatalf("decode downloaded audio: %v", err)
			}
			if uploaded.SampleRate == nil || downloadedMeta.SampleRate != uint32(*uploaded.SampleRate) {
				t.Fatalf("downloaded sample rate = %d, want %v", downloadedMeta.SampleRate, uploaded.SampleRate)
			}

			// flac 编码器写 STREAMINFO 时 NSamples=0，DecodeMeta 时长恒为 0，
			// 因此时长用全量解码后的样本数验证，同时确认音频内容完整。
			full, err := audioio.NewDecoder().Decode(bytes.NewReader(downloadResponse.Body.Bytes()))
			if err != nil {
				t.Fatalf("decode downloaded audio data: %v", err)
			}
			if uploaded.Duration != nil && downloadedMeta.SampleRate > 0 {
				got := float32(len(full)) / float32(downloadedMeta.SampleRate)
				if delta := got - *uploaded.Duration; delta < -0.5 || delta > 0.5 {
					t.Fatalf("downloaded audio duration = %.2fs, want %.2fs (±0.5s)", got, *uploaded.Duration)
				}
			}

			deleteRequest := httptest.NewRequest(http.MethodDelete, "/audio/"+*uploaded.ID, nil)
			deleteResponse := httptest.NewRecorder()
			router.ServeHTTP(deleteResponse, deleteRequest)
			if deleteResponse.Code != http.StatusNoContent {
				t.Fatalf("delete status = %d, want %d", deleteResponse.Code, http.StatusNoContent)
			}
		})
	}
}

func multipartFixture(t *testing.T, name string, content []byte) (*bytes.Reader, string) {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(body.Bytes()), writer.FormDataContentType()
}
