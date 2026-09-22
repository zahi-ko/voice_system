package handlers

import (
	"github.com/labstack/echo/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"voice_system/internal/transport/http/generated"
)

func (h *Handler) ApplyEffect(ctx *echo.Context, audioID openapi_types.UUID, params generated.ApplyEffectParams) error {
	return notImplemented()
}

func (h *Handler) ApplyEffectChain(ctx *echo.Context, audioID openapi_types.UUID, params generated.ApplyEffectChainParams) error {
	return notImplemented()
}

func (h *Handler) GetEffectHistory(ctx *echo.Context, audioID openapi_types.UUID) error {
	return notImplemented()
}

func (h *Handler) ListEffects(ctx *echo.Context, audioID openapi_types.UUID) error {
	return notImplemented()
}

func (h *Handler) UndoEffect(ctx *echo.Context, audioID openapi_types.UUID) error {
	return notImplemented()
}
