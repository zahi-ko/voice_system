package main

import (
	"log"

	"voice_system/internal/adapters/audioio"
	"voice_system/internal/adapters/memory"
	"voice_system/internal/application"
	"voice_system/internal/transport/http/handlers"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()
	store, err := memory.NewAudioStore("./tmp")
	if err != nil {
		log.Fatal(err)
	}
	service := application.NewAudioService(audioio.NewWAVDecoder(), store)
	handler := handlers.NewHandler(service)
	handlers.Register(e, handler)

	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
