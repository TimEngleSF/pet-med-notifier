package main

import (
	"log"
	"net/http"
	"os"

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

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!\n")
	})

	e.Logger.Fatal(e.Start(":42069"))
}
