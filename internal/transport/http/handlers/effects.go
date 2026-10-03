package handlers

import (
	"encoding/json"
	"errors"
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
		// 效果名/参数问题属于调用方错误，回 400 并带出具体原因，其余按 500 处理
		if isClientEffectError(err) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
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
		if isClientEffectError(err) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to apply effect chain")
	}

	response := transporthttp.EffectChainResponseToAPI(aud, applied)

	middleware.SetDetail(ctx, "chain=%d save=%t applied=%s", len(payload), save, strings.TrimSuffix(applied, "->"))

	return ctx.JSON(http.StatusOK, response)
}

// isClientEffectError 判断效果失败是否由请求内容引起（效果名不存在、参数缺失或越界）。
func isClientEffectError(err error) bool {
	return errors.Is(err, effects.ErrUnknownEffect) || errors.Is(err, effects.ErrInvalidParameters)
}

func (h *Handler) ListEffects(ctx *echo.Context) error {
	if h.effectService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "effect service unavailable")
	}

	effectsList := h.effectService.ListAvailable()
	middleware.SetDetail(ctx, "count=%d", len(effectsList))
	return ctx.JSON(http.StatusOK, effectsList)
}
