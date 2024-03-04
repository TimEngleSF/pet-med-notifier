package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/TimEngleSF/pet-med-notifier/repository"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client    *mongo.Client
	MedDb     *mongo.Database
	DummyMeds = repository.Medicines{
		{
			Id:         primitive.NewObjectID(),
			Name:       "Test1",
			TimeToTake: &repository.TimeKey{Hour: 00, Min: 00},
			Taken:      false,
			Due:        true,
		},
		{
			Id:         primitive.NewObjectID(),
			Name:       "Test2",
			TimeToTake: &repository.TimeKey{Hour: 16, Min: 00},
			Taken:      false,
			Due:        false,
		},
		{
			Id:         primitive.NewObjectID(),
			Name:       "Test3",
			TimeToTake: &repository.TimeKey{Hour: 00, Min: 00},
			Taken:      true,
			Due:        true,
		},
		{
			Id:         primitive.NewObjectID(),
			Name:       "Test2",
			TimeToTake: &repository.TimeKey{Hour: 16, Min: 00},
			Taken:      false,
			Due:        false,
		},
		{
			Id:         primitive.NewObjectID(),
			Name:       "Test3",
			TimeToTake: &repository.TimeKey{Hour: 00, Min: 00},
			Taken:      true,
			Due:        true,
		},
	}
)

func main() {
	godotenv.Load()
	// Set
	URI := os.Getenv("URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(URI))
	MedDb = client.Database("lily-med")
	// Run your server.
	if err := runServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}
}
