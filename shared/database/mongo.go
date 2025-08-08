package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// MongoConfig contains MongoDB connection configuration
type MongoConfig struct {
	URI      string `env:"MONGODB_URI"`
	Database string `env:"MONGODB_DB"`
}

// NewMongo creates a new MongoDB configuration from environment variables
func NewMongoConfig() *MongoConfig {
	cfg, err := env.ParseAs[MongoConfig]()
	if err != nil {
		log.Fatalf("Failed to parse env: %v", err)
	}

	return &cfg
}

func NewMongoClient(ctx context.Context, cfg *MongoConfig) (*mongo.Client, error) {
	if cfg.URI == "" {
		return nil, fmt.Errorf("mongo URI is not provided")
	}
	if cfg.Database == "" {
		return nil, fmt.Errorf("mongo database is not provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	log.Printf("🎉 Successfully connected to MongoDB at %s", cfg.URI)

	return client, nil
}

func GetDatabase(client *mongo.Client, cfg *MongoConfig) *mongo.Database {
	return client.Database(cfg.Database)
}
