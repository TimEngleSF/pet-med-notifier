package main

import (
	"context"
	db "lily-med/DB"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var isTest bool = false

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	env := os.Getenv("GO_ENV")
	if env == "test" {
		isTest = true
	}
	// TODO: connect to db and handle errors, then begin writing tests

	dbInstance, err := db.GetInstance(context.Background(), isTest)
	if err != nil {
		log.Fatalf("Error initializing database connection: %v\n", err)
	}

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			dbInstance.PingDatabase()
		}
	}()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!\n")
	})

	e.Logger.Fatal(e.Start(":42069"))
}
