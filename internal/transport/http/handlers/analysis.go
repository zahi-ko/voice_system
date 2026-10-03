package handlers

import (
	"net/http"
	transporthttp "voice_system/internal/transport/http"
	"voice_system/internal/transport/http/generated"
	"voice_system/internal/transport/http/middleware"

	"github.com/labstack/echo/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handler) GetSpectrum(ctx *echo.Context, audioID openapi_types.UUID, params generated.GetSpectrumParams) error {
	if h.analysisService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "analysis service not available")
	}

	id := transporthttp.APItoUUID(audioID)
	nfft := transporthttp.APItoNFFT(params)
	spectrum, err := h.analysisService.GetSpectrum(ctx.Request().Context(), id, nfft)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The request was invalid or cannot be served.")
	}

	middleware.SetDetail(ctx, "id=%s nfft=%d bins=%d", id, nfft, len(spectrum.Data))
	return ctx.JSON(http.StatusOK, spectrum)
}

func (h *Handler) GetStats(ctx *echo.Context, audioID openapi_types.UUID) error {
	if h.analysisService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "analysis service not available")
	}

	id := transporthttp.APItoUUID(audioID)
	stats, err := h.analysisService.GetAudioStats(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The request was invalid or cannot be served.")
	}

	middleware.SetDetail(ctx, "id=%s rms=%.3f peak=%.3f", id, stats.RMS, stats.Peak)
	return ctx.JSON(http.StatusOK, stats)
}

func (h *Handler) GetWaveform(ctx *echo.Context, audioID openapi_types.UUID, params generated.GetWaveformParams) error {
	if h.analysisService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "analysis service not available")
	}

	id := transporthttp.APItoUUID(audioID)
	points := transporthttp.APItoPoint(params)
	waveform, err := h.analysisService.GetWaveForm(ctx.Request().Context(), id, points)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The request was invalid or cannot be served.")
	}

	middleware.SetDetail(ctx, "id=%s points=%d", id, points)
	return ctx.JSON(http.StatusOK, waveform)
}
