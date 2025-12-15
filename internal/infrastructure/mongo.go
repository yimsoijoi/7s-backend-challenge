package infrastructure

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

func NewMongoDatabase(cfg MongoConfig) (*mongo.Client, *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		log.Fatalf("mongo connect failed: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping failed: %v", err)
	}

	log.Println("MongoDB connected")

	return client, client.Database(cfg.Database)
}
