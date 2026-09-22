package handlers

import (
	"voice_system/internal/transport/http/generated"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetSpectrogram(ctx *echo.Context, audioID string, params generated.GetSpectrogramParams) error {
	return notImplemented()
}

func (h *Handler) GetSpectrum(ctx *echo.Context, audioID string, params generated.GetSpectrumParams) error {
	return notImplemented()
}

func (h *Handler) GetStats(ctx *echo.Context, audioID string) error {
	return notImplemented()
}

func (h *Handler) GetWaveform(ctx *echo.Context, audioID string, params generated.GetWaveformParams) error {
	return notImplemented()
}
