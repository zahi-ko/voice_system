package handlers

import (
	"net/http"
	"voice_system/internal/domain/effects"
	"voice_system/internal/transport/http/generated"

	"github.com/labstack/echo/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"

	transporthttp "voice_system/internal/transport/http"
)

func (h *Handler) ApplyEffect(ctx *echo.Context, audioID openapi_types.UUID, params generated.ApplyEffectParams) error {
	if h.effectService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "effect service unavailable")
	}

	var payload effects.Payload
	if err := ctx.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	id := transporthttp.APItoUUID(audioID)
	save := transporthttp.APItoSaveAsNew(params)

	aud, err := h.effectService.Apply(ctx.Request().Context(), id, save, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to apply effect")
	}

	audioMeta := transporthttp.AudioToAPI(aud)

	return ctx.JSON(http.StatusOK, audioMeta)
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

// func bindEffect(ctx *echo.Context, )
