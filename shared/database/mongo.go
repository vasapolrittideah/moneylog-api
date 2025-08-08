package database

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// MongoConfig contains MongoDB connection configuration.
type MongoConfig struct {
	URI      string `env:"MONGODB_URI"`
	Database string `env:"MONGODB_DB"`
}

const (
	DefaultConnectionTimeout = 20 * time.Second
)

// NewMongoConfig creates a new MongoDB configuration from environment variables.
func NewMongoConfig() *MongoConfig {
	cfg, err := env.ParseAs[MongoConfig]()
	if err != nil {
		slog.With("error", err).Error("Failed to parse env")
		panic(err)
	}

	return &cfg
}

// Validate checks if the MongoDB configuration is valid.
func (c *MongoConfig) Validate() error {
	if c.URI == "" {
		return errors.New("mongo URI is not provided")
	}
	if c.Database == "" {
		return errors.New("mongo database is not provided")
	}
	return nil
}

// MongoDB implements the Database interface for MongoDB.
type MongoDB struct {
	config   *MongoConfig
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDB creates a new MongoDB instance.
func NewMongoDB(config *MongoConfig) *MongoDB {
	return &MongoDB{
		config: config,
	}
}

// Connect establishes a connection to MongoDB.
func (m *MongoDB) Connect(ctx context.Context) error {
	if err := m.config.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, DefaultConnectionTimeout)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(m.config.URI))
	if err != nil {
		return err
	}

	m.client = client
	m.database = client.Database(m.config.Database)

	if pingErr := m.Ping(ctx); pingErr != nil {
		return pingErr
	}

	logger := slog.Default()
	logger.InfoContext(ctx, "Connected to MongoDB", "uri", m.config.URI)
	return nil
}

// Disconnect closes the MongoDB connection.
func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.client == nil {
		return nil
	}
	return m.client.Disconnect(ctx)
}

// Ping checks if the MongoDB connection is alive.
func (m *MongoDB) Ping(ctx context.Context) error {
	if m.client == nil {
		return errors.New("client is not connected")
	}
	return m.client.Ping(ctx, readpref.Primary())
}
