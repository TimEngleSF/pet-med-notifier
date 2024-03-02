package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	gowebly "github.com/gowebly/helpers"
)

var client *mongo.Client
var MedDb *mongo.Database

// runServer runs a new HTTP server with the loaded environment variables.
func runServer() error {
	URI := os.Getenv("URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(URI))
	MedDb = client.Database("lily-med")
	fmt.Println(MedDb.Name())
	// Validate environment variables.
	port, err := strconv.Atoi(gowebly.Getenv("BACKEND_PORT", "42069"))
	if err != nil {
		return err
	}

	// Create a new Echo server.
	echo := echo.New()

	// Add Echo middlewares.
	echo.Use(middleware.Logger())

	// Handle static files.
	echo.Static("/static", "./static")

	// Handle index page view.
	echo.GET("/", indexViewHandler)

	// Handle API endpoints.
	echo.GET("/api/hello-world", showContentAPIHandler)

	// Create a new server instance with options from environment variables.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      echo, // handle all Echo routes
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Send log message.
	slog.Info("Starting server...", "port", port)

	return server.ListenAndServe()
}
