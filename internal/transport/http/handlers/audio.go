package handlers

import (
	"github.com/labstack/echo/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handler) ListAudio(ctx *echo.Context) error {
	return notImplemented()
}

func (h *Handler) UploadAudio(ctx *echo.Context) error {
	return notImplemented()
}

func (h *Handler) RemoveAudio(ctx *echo.Context, audioID openapi_types.UUID) error {
	return notImplemented()
}

func (h *Handler) DownloadAudio(ctx *echo.Context, audioID openapi_types.UUID) error {
	return notImplemented()
}
