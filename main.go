package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// Run your server.
	if err := runServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}
}
