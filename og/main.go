package main

import (
	"context"
	db "lily-med/DB"
	"lily-med/controllers"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

var isTest bool = false

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	dbInstance, err := db.GetInstance(context.Background())
	if err != nil {
		log.Fatalf("Error initializing database connection: %v\n", err)
	}

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			dbInstance.PingDatabase()
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!\n"))
	})

	mux.HandleFunc("/medicine", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.CreateMedicineHandler(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	server := &http.Server{Addr: ":42069", Handler: mux}
	log.Println("Starting server on 42069")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
