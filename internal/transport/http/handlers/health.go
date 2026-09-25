package handlers

import (
	"net/http"
	"voice_system/internal/transport/http/middleware"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetHealth(ctx *echo.Context) error {
	if h.audioService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "audio service unavailable")
	}

	middleware.SetDetail(ctx, "status=ok")
	return ctx.NoContent(http.StatusOK)
}
