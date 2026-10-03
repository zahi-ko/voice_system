package handlers

import (
	"net/http"

	"voice_system/internal/application"
	"voice_system/internal/transport/http/generated"

	"github.com/labstack/echo/v5"
)

// Handler contains the HTTP handlers for the API.
type Handler struct {
	audioService    *application.AudioService
	effectService   *application.EffectService
	analysisService *application.AnalysisService
}

var _ generated.ServerInterface = (*Handler)(nil)

// NewHandler creates an API handler skeleton.
func NewHandler(audioService *application.AudioService, effectService *application.EffectService, analysisService *application.AnalysisService) *Handler {
	return &Handler{
		audioService:    audioService,
		effectService:   effectService,
		analysisService: analysisService,
	}
}

// Register attaches all generated API routes to the supplied Echo router.
func Register(router generated.EchoRouter, handler *Handler) {
	generated.RegisterHandlers(router, handler)
}

func notImplemented() error {
	return echo.NewHTTPError(http.StatusNotImplemented, "handler not implemented")
}
