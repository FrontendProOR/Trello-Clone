package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() *mongo.Database {
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:pass@localhost:27017"))
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:pass@mongodb:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("mongoDB")
	return DB
}
