package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"voice_system/internal/transport/http/generated"
)

// Handler contains the HTTP handlers for the API.
type Handler struct{}

var _ generated.ServerInterface = (*Handler)(nil)

// NewHandler creates an API handler skeleton.
func NewHandler() *Handler {
	return &Handler{}
}

// Register attaches all generated API routes to the supplied Echo router.
func Register(router generated.EchoRouter, handler *Handler) {
	generated.RegisterHandlers(router, handler)
}

func notImplemented() error {
	return echo.NewHTTPError(http.StatusNotImplemented, "handler not implemented")
}
