package database

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var URI = os.Getenv("LILY_MONGO_URI")
var (
	instance *Db
	db       *mongo.Database
	client   *mongo.Client
	dbName   string
)

type Db struct {
	Db              *mongo.Database
	Client          *mongo.Client
	CollectionNames []string
}

func InitDatabase(c context.Context) (*Db, error) {
	if URI == "" {
		return nil, errors.New("MONGO_URI environment variable must be set")
	}
	isTest := false
	env := os.Getenv("GO_ENV")
	if env == "test" {
		isTest = true
	}

	dbName = "lily-med"
	if isTest {
		dbName = "lily-med-test"
	}

	client := ConnectClient(c)

	db = ConnectDatabase(client, dbName)

	return &Db{
		Client:          client,
		Db:              db,
		CollectionNames: []string{"medications", "history"},
	}, nil
}

func ConnectClient(c context.Context) *mongo.Client {
	c, cancel := context.WithCancel(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Client(c, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Connection to MongoDB Client %v Failed: %v\n", uri, err)
	}

	return client
}

func ConnectDatabase(c *mongo.Client, n string) *mongo.Database {
	return c.Database(n)
}

func GetInstance(c context.Context) (*Db, error) {
	instance, err := InitDatabase(c)
	return instance, err
}

func (d Db) CloseConnection() {
	err := d.Client.Disconnect(context.TODO())
	if err != nil {
		log.Printf("Error closing mongo client: %v\n", err)
	}
}

// func GetInstance (c *context.Context) (*Db, error) {
// 	if
// }
