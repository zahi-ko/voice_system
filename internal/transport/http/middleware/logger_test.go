package middleware

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
)

// captureLog redirects the std logger output while fn runs.
func captureLog(t *testing.T, fn func()) string {
	t.Helper()

	var buf bytes.Buffer
	flags := log.Flags()
	out := log.Writer()
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		log.SetOutput(out)
		log.SetFlags(flags)
	}()

	fn()
	return buf.String()
}

func TestHumanSize(t *testing.T) {
	cases := []struct {
		n    int64
		want string
	}{
		{0, "0B"},
		{512, "512B"},
		{1 << 10, "1.0KB"},
		{1536 << 10, "1.5MB"},
	}
	for _, c := range cases {
		if got := HumanSize(c.n); got != c.want {
			t.Errorf("HumanSize(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

func TestStatusColor(t *testing.T) {
	cases := map[int]string{
		200: ansiGreen,
		201: ansiGreen,
		404: ansiYel,
		422: ansiYel,
		500: ansiRed,
	}
	for code, color := range cases {
		if got := statusColor(code); got != color {
			t.Errorf("statusColor(%d) color = %q, want %q", code, got, color)
		}
	}
}

func newLoggingRouter(t *testing.T) *echo.Echo {
	t.Helper()
	router := echo.New()
	router.Use(RequestLogger())
	return router
}

func TestRequestLoggerSuccessWithDetail(t *testing.T) {
	router := newLoggingRouter(t)
	router.POST("/audio", func(c *echo.Context) error {
		SetDetail(c, "name=song fmt=wav size=1.0KB")
		return c.JSON(http.StatusCreated, map[string]string{"ok": "true"})
	})

	logLine := captureLog(t, func() {
		req := httptest.NewRequest(http.MethodPost, "/audio", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
		}
	})

	for _, want := range []string{"[http]", "POST", "/audio", "201", "name=song", "fmt=wav", "size=1.0KB"} {
		if !strings.Contains(logLine, want) {
			t.Errorf("log line %q missing %q", stripANSI(logLine), want)
		}
	}
}

func TestRequestLoggerHTTPErrorDerivation(t *testing.T) {
	router := newLoggingRouter(t)
	router.GET("/audio/:id", func(c *echo.Context) error {
		// 不写响应，仅返回错误：状态码应从 echo.HTTPError 推导。
		return echo.NewHTTPError(http.StatusNotFound, "audio not found")
	})

	logLine := captureLog(t, func() {
		req := httptest.NewRequest(http.MethodGet, "/audio/abc", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	})

	for _, want := range []string{"GET", "/audio/abc", "404", "err=audio not found"} {
		if !strings.Contains(logLine, want) {
			t.Errorf("log line %q missing %q", stripANSI(logLine), want)
		}
	}
}

func TestRequestLoggerInternalErrorDerivation(t *testing.T) {
	router := newLoggingRouter(t)
	router.GET("/boom", func(c *echo.Context) error {
		return errors.New("disk on fire")
	})

	logLine := captureLog(t, func() {
		req := httptest.NewRequest(http.MethodGet, "/boom", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	})

	for _, want := range []string{"500", "err=disk on fire"} {
		if !strings.Contains(logLine, want) {
			t.Errorf("log line %q missing %q", stripANSI(logLine), want)
		}
	}
}

func stripANSI(s string) string {
	var b strings.Builder
	inSeq := false
	for _, r := range s {
		switch {
		case r == '\033':
			inSeq = true
		case inSeq && r == 'm':
			inSeq = false
		case !inSeq:
			b.WriteRune(r)
		}
	}
	return b.String()
}
