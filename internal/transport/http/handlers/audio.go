package handlers

import (
	"net/http"

	transporthttp "voice_system/internal/transport/http"

	"github.com/labstack/echo/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handler) ListAudio(ctx *echo.Context) error {
	if h.audioService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "audio service unavailable")
	}

	audioList, err := h.audioService.List(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The request was invalid or cannot be served.")
	}

	return ctx.JSON(http.StatusOK, transporthttp.AudioListToAPI(audioList))
}

func (h *Handler) UploadAudio(ctx *echo.Context) error {
	if h.audioService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "audio service unavailable")
	}

	file, header, err := ctx.Request().FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "file is required")
	}
	defer file.Close()

	item, err := h.audioService.Upload(ctx.Request().Context(), header.Filename, file)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, transporthttp.AudioToAPI(item))
}

func (h *Handler) RemoveAudio(ctx *echo.Context, audioID openapi_types.UUID) error {
	return notImplemented()
}

func (h *Handler) DownloadAudio(ctx *echo.Context, audioID openapi_types.UUID) error {
	return notImplemented()
}
