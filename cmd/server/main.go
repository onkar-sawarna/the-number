package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thenumber/app/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "47321"
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "web/static")

	handlers.New().Register(e)

	addr := "0.0.0.0:" + port
	log.Printf("the number listening on %s", addr)
	e.Logger.Fatal(e.Start(addr))
}
