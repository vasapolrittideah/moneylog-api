package discovery

import (
	"google.golang.org/grpc"
)

// Registry defines the interface for service discovery.
type Registry interface {
	// Register registers a service instance with the registry
	Register(instanceID string, serviceName string, hostPort string) error

	// Deregister removes a service instance from the registry
	Deregister(instanceID string, serviceName string) error

	// Connect establishes a connection to a service
	Connect(hostPort string, serviceName string) (*grpc.ClientConn, error)

	// Close gracefully shuts down the registry
	Close() error
}
