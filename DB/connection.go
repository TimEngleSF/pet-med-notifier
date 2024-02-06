package db

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	instance *Db
	db       *mongo.Database
	client   *mongo.Client
	uri      string
	dbName   string
	once     sync.Once
)

type Db struct {
	Db     *mongo.Database
	Client *mongo.Client
	Config DbConf
}

type DbConf struct {
	Uri    string
	DbName string
}

func initDatabase(ctx context.Context, isTestEnv bool) (*Db, error) {
	uri = os.Getenv("MONGO_URI")
	if uri == "" {
		return nil, errors.New("MONGO_URI environment variable must be set")
	}

	dbName = "lily-med"
	if isTestEnv {
		dbName = "lily-med-test"
	}

	client = ConnectClient(ctx, uri)

	db = ConnectDatabase(client, dbName)

	return &Db{
		Client: client,
		Db:     db,
		Config: DbConf{
			Uri:    uri,
			DbName: dbName,
		},
	}, nil
}

func ConnectClient(ctx context.Context, uri string) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}

func ConnectDatabase(c *mongo.Client, n string) *mongo.Database {
	return c.Database(n)
}

func GetInstance(ctx context.Context, isTestEnv bool) (*Db, error) {
	var err error
	once.Do(func() {
		instance, err = initDatabase(ctx, isTestEnv)
	})
	return instance, err
}

func (d Db) CloseConnection() {
	err := d.Client.Disconnect(context.TODO())
	if err != nil {
		log.Printf("Error closing mongo client: %v", err)
	}
}
