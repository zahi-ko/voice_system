package handlers

import (
	"net/http"
	"os"

	"voice_system/internal/domain/audio"
	"voice_system/internal/transport/http/middleware"
	transporthttp "voice_system/internal/transport/http"

	"github.com/labstack/echo/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

var mimeMp = map[audio.Format]string{
	"mp3":  "audio/mpeg",
	"wav":  "audio/wav",
	"ogg":  "audio/ogg",
	"flac": "audio/flac",
	"aiff": "audio/aiff",
}

func (h *Handler) ListAudio(ctx *echo.Context) error {
	if h.audioService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "audio service unavailable")
	}

	audioList, err := h.audioService.List(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The request was invalid or cannot be served.")
	}

	middleware.SetDetail(ctx, "count=%d", len(audioList.Items))
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

	middleware.SetDetail(ctx, "name=%s fmt=%s size=%s dur=%.1fs",
		item.Name, item.Meta.Format, middleware.HumanSize(header.Size), item.Meta.Duration)
	return ctx.JSON(http.StatusCreated, transporthttp.AudioToAPI(item))
}

func (h *Handler) RemoveAudio(ctx *echo.Context, audioID openapi_types.UUID) error {
	if h.audioService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "audio service unavailable")
	}

	id := transporthttp.APItoUUID(audioID)

	err := h.audioService.Delete(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "The request was invalid or cannot be serverd")
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *Handler) DownloadAudio(ctx *echo.Context, audioID openapi_types.UUID) error {
	if h.audioService == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "audio service unavailable")
	}

	id := transporthttp.APItoUUID(audioID)
	path, err := h.audioService.Download(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "the requested resource is not found")
	}

	a, err := h.audioService.Lookup(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "audio file not found")
	}

	mime, ok := mimeMp[a.Meta.Format]
	if !ok {
		mime = "application/octet-stream"
	}

	file, err := os.Open(path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "file not found")
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "stat failed")
	}

	name := a.Name + "." + string(a.Meta.Format)
	modtime := info.ModTime()

	middleware.SetDetail(ctx, "name=%s size=%s", name, middleware.HumanSize(info.Size()))

	ctx.Response().Header().Set(echo.HeaderContentType, mime)
	// ctx.Response().Header().Set(echo.HeaderContentDisposition,
	// 	fmt.Sprintf(`attachment; filename="%s"`, name))

	http.ServeContent(ctx.Response(), ctx.Request(), name, modtime, file)
	return nil
}
