package main

import (
	"io"
	"log"
	"log/slog"

	"voice_system/internal/adapters/audioio"
	"voice_system/internal/adapters/memory"
	"voice_system/internal/application"
	"voice_system/internal/transport/http/handlers"
	"voice_system/internal/transport/http/middleware"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()
	// 框架自带日志丢弃，请求日志由 middleware.RequestLogger 接管
	e.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	e.Use(middleware.RequestLogger())

	store, err := memory.NewAudioStore()
	if err != nil {
		log.Fatal(err)
	}
	service := application.NewAudioService(audioio.NewDecoder(), store)
	handler := handlers.NewHandler(service)
	handlers.Register(e, handler)
	e.Static("/", "web/dist")

	log.Print("listening on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
