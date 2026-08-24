package main

import (
	"log"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thenumber/app/internal/handlers"
	"github.com/thenumber/app/internal/models"
	appsession "github.com/thenumber/app/internal/session"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "47321"
	}
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "data/thenumber.db"
	}

	db, err := models.Open(dbPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(appsession.NewStore()))
	e.Static("/static", "web/static")

	handlers.New(db).Register(e)

	addr := "0.0.0.0:" + port
	log.Printf("the number listening on %s", addr)
	e.Logger.Fatal(e.Start(addr))
}
