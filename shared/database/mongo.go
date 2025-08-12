package database

import (
	"context"
	"errors"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
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
func NewMongoConfig(logger *zerolog.Logger) *MongoConfig {
	cfg, err := env.ParseAs[MongoConfig]()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse env")
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
	logger   *zerolog.Logger
}

// NewMongoDB creates a new MongoDB instance.
func NewMongoDB(config *MongoConfig, logger *zerolog.Logger) *MongoDB {
	return &MongoDB{
		config: config,
		logger: logger,
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

	if err := m.client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	m.logger.Info().Str("uri", m.config.URI).Msg("Connected to MongoDB")
	return nil
}

// Disconnect closes the MongoDB connection.
func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.client == nil {
		return nil
	}

	m.logger.Info().Msg("Disconnected from MongoDB")
	return m.client.Disconnect(ctx)
}

// GetDatabase returns the MongoDB database.
func (m *MongoDB) GetDatabase() *mongo.Database {
	return m.database
}
