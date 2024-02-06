package helpers

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

func LoadEnv() {
	once.Do(func() {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("no .env file found")
			panic(1)
		}
	})
}

func GetURIString() (string, error) {
	LoadEnv()
	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		log.Println("MONGO_URI variable does not exist")
		return "", nil
	}

	return uri, nil
}
