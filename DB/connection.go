package db

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"time"

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

func initDatabase(ctx context.Context) (*Db, error) {
	uri = os.Getenv("MONGO_URI")
	if uri == "" {
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

	client, err := ConnectClient(ctx, uri)
	if err != nil {
		return nil, err
	}

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

func ConnectClient(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Printf("Error connecting to client: %v", err)
		return nil, err
	}
	return client, nil
}

func ConnectDatabase(c *mongo.Client, n string) *mongo.Database {
	return c.Database(n)
}

func GetInstance(ctx context.Context) (*Db, error) {
	var err error
	once.Do(func() {
		instance, err = initDatabase(ctx)
	})
	return instance, err
}

func (d Db) CloseConnection() {
	err := d.Client.Disconnect(context.TODO())
	if err != nil {
		log.Printf("Error closing mongo client: %v", err)
	}
}

func (d *Db) PingDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := d.Client.Ping(ctx, nil); err != nil {
		log.Printf("Error pinging database: %v\n", err)
		client, err = ConnectClient(ctx, d.Config.Uri)
		if err != nil {
			log.Printf("Error reconnecting to mongo client after failed database ping: %v\n", err)
		}

		// TODO: What should happen if the client doesnt connect again? Might be worth pinging here and if it fails send an email out to admin
	} else {
		log.Println("Successfully pinged database.")
	}
}
