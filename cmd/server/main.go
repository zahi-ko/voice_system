package main

import (
	"log"

	"voice_system/internal/transport/http/handlers"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()
	handler := handlers.NewHandler()
	handlers.Register(e, handler)

	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
