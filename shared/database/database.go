package database

import (
	"context"
)

// Database represents a generic database interface that can be implemented
// by different database providers (MongoDB, PostgreSQL, etc.)
type Database interface {
	// Connect establishes a connection to the database
	Connect(ctx context.Context) error

	// Disconnect closes the database connection
	Disconnect(ctx context.Context) error

	// Ping checks if the database connection is alive
	Ping(ctx context.Context) error
}

// Config represents generic database configuration.
type Config interface {
	Validate() error
}
