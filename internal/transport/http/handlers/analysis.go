package handlers

import (
	"net/http"
	transporthttp "voice_system/internal/transport/http"
	"voice_system/internal/transport/http/generated"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetSpectrum(ctx *echo.Context, audioID string, params generated.GetSpectrumParams) error {
	if h.analysisService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "analysis service not available")
	}

	nfft := transporthttp.APItoNFFT(params)
	spectrum, err := h.analysisService.GetSpectrum(ctx.Request().Context(), audioID, nfft)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The request was invalid or cannot be served.")
	}

	return ctx.JSON(http.StatusOK, spectrum)
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
