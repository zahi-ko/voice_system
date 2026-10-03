package handlers

import (
	"net/http"
	"voice_system/internal/transport/http/generated"
	transporthttp "voice_system/internal/transport/http"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetSpectrogram(ctx *echo.Context, audioID string, params generated.GetSpectrogramParams) error {
	return notImplemented()
}

func (h *Handler) GetSpectrum(ctx *echo.Context, audioID string, params generated.GetSpectrumParams) error {
	return notImplemented()
}

func (h *Handler) GetStats(ctx *echo.Context, audioID string) error {
	if h.analysisService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "analysis service not available")
	}

	stats, err := h.analysisService.GetAudioStats(ctx.Request().Context(), audioID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The request was invalid or cannot be served.")
	}

	return ctx.JSON(http.StatusOK, stats)
}

func (h *Handler) GetWaveform(ctx *echo.Context, audioID string, params generated.GetWaveformParams) error {
	if h.analysisService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "analysis service not available")
	}

	points := transporthttp.APItoPoint(params)
	waveform, err := h.analysisService.GetWaveForm(ctx.Request().Context(), audioID, points)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The request was invalid or cannot be served.")
	}

	return ctx.JSON(http.StatusOK, waveform)
}
