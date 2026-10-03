package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"voice_system/internal/domain/effects"
	"voice_system/internal/transport/http/generated"
	"voice_system/internal/transport/http/middleware"

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

	effectName, _ := payload["name"].(string)
	middleware.SetDetail(ctx, "effect=%s save=%t", effectName, save)

	audioMeta := transporthttp.AudioToAPI(aud)

	return ctx.JSON(http.StatusOK, audioMeta)
}

func (h *Handler) ApplyEffectChain(ctx *echo.Context, audioID openapi_types.UUID, params generated.ApplyEffectChainParams) error {
	if h.effectService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "effect service unavailable")
	}

	var payload []effects.Payload
	if err := json.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	if len(payload) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "effect chain cannot be empty")
	}

	id := transporthttp.APItoUUID(audioID)
	save := transporthttp.APItoSaveAsNewChain(params)

	aud, applied, err := h.effectService.ApplyChain(ctx.Request().Context(), id, save, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to apply effect chain")
	}

	response := transporthttp.EffectChainResponseToAPI(aud, applied)

	middleware.SetDetail(ctx, "chain=%d save=%t applied=%s", len(payload), save, strings.TrimSuffix(applied, "->"))

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handler) ListEffects(ctx *echo.Context) error {
	if h.effectService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "effect service unavailable")
	}

	effectsList := h.effectService.ListAvailable()
	middleware.SetDetail(ctx, "count=%d", len(effectsList))
	return ctx.JSON(http.StatusOK, effectsList)
}
