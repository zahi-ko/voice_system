package handlers

import (
	"net/http"
	"voice_system/internal/domain/effects"
	"voice_system/internal/transport/http/generated"

	"github.com/labstack/echo/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handler) ApplyEffect(ctx *echo.Context, audioID openapi_types.UUID, params generated.ApplyEffectParams) error {
	if h.effectService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "effect service unavailable")
	}

	// id := transporthttp.APItoUUID(audioID)
	// save := transporthttp.APItoSaveAsNew(params)

	var effect effects.Effect
	if err := ctx.Bind(&effect); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid effect parameters")
	}

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
