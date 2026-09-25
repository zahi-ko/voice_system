// Package middleware provides HTTP middleware for the transport layer.
package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

// ANSI color codes for the request logger.
const (
	ansiReset = "\033[0m"
	ansiDim   = "\033[90m"
	ansiCyan  = "\033[36m"
	ansiBlue  = "\033[34m"
	ansiGreen = "\033[32m"
	ansiYel   = "\033[33m"
	ansiRed   = "\033[31m"
	ansiMag   = "\033[35m"
)

func paint(color, s string) string { return color + s + ansiReset }

func statusColor(code int) string {
	switch {
	case code >= 500:
		return ansiRed
	case code >= 400:
		return ansiYel
	default:
		return ansiGreen
	}
}

// detailKey is the context key handlers use to attach business detail
// that gets appended to the request log line.
const detailKey = "middleware.log_detail"

// SetDetail attaches a business-detail suffix to the current request's
// log line, e.g. middleware.SetDetail(c, "name=%s size=%s", name, size).
// Call it before returning; the request logger prints it after the
// status and latency fields.
func SetDetail(c *echo.Context, format string, args ...any) {
	c.Set(detailKey, fmt.Sprintf(format, args...))
}

// RequestLogger returns a middleware that prints one concise, colorized
// line per request:
//
//	[http] POST /audio 201 45.2ms name=song.mp3 fmt=mp3 size=1.2MB
//
// Status is taken from the written response when available, otherwise
// derived from the returned error (echo.HTTPError code or 500).
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			req := c.Request()

			err := next(c)

			// Status: read the written response when committed, otherwise
			// derive from the returned error (echo.HTTPError code or 500).
			status := 0
			if resp, uerr := echo.UnwrapResponse(c.Response()); uerr == nil {
				if resp.Committed {
					status = resp.Status
				} else {
					status = http.StatusOK // implicit default if no error
				}
			}
			if err != nil && (status == 0 || status == http.StatusOK) {
				var he *echo.HTTPError
				switch {
				case errors.As(err, &he):
					status = he.Code
				default:
					status = http.StatusInternalServerError
				}
			}

			line := fmt.Sprintf("%s %s %s %s %s",
				paint(ansiDim, "[http]"),
				paint(ansiCyan, fmt.Sprintf("%-7s", req.Method)),
				paint(ansiBlue, req.URL.Path),
				paint(statusColor(status), fmt.Sprintf("%3d", status)),
				paint(ansiMag, time.Since(start).Round(100*time.Microsecond).String()),
			)

			// Business detail attached by the handler (if any).
			if d, ok := c.Get(detailKey).(string); ok && d != "" {
				line += " " + paint(ansiDim, d)
			}

			// Error message for failed requests.
			if err != nil {
				var he *echo.HTTPError
				if errors.As(err, &he) {
					line += " " + paint(ansiRed, "err="+fmt.Sprint(he.Message))
				} else {
					line += " " + paint(ansiRed, "err="+err.Error())
				}
			}

			log.Print(line)
			return err
		}
	}
}

// HumanSize formats a byte count as a compact human-readable string.
func HumanSize(n int64) string {
	const kb, mb = 1 << 10, 1 << 20
	switch {
	case n >= mb:
		return fmt.Sprintf("%.1fMB", float64(n)/mb)
	case n >= kb:
		return fmt.Sprintf("%.1fKB", float64(n)/kb)
	default:
		return fmt.Sprintf("%dB", n)
	}
}
