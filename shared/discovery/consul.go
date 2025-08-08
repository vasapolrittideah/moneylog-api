package discovery

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	// Required for consul:// resolver to work with gRPC.
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// HealthCheckConfig defines health check parameters.
type HealthCheckConfig struct {
	Interval        string
	Timeout         string
	DeregisterAfter string
}

// DefaultHealthCheckConfig returns sensible defaults.
func DefaultHealthCheckConfig() *HealthCheckConfig {
	return &HealthCheckConfig{
		Interval:        "10s",
		Timeout:         "5s",
		DeregisterAfter: "1m",
	}
}

// ConsulRegistry defines a consul based service registry.
type ConsulRegistry struct {
	client      *consulapi.Client
	healthCheck *HealthCheckConfig
	logger      *slog.Logger
}

// NewConsulRegistry creates a new Consul registry with default health check settings.
func NewConsulRegistry(addr string) (*ConsulRegistry, error) {
	cfg := consulapi.DefaultConfig()
	cfg.Address = addr

	client, err := consulapi.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ConsulRegistry{
		client:      client,
		healthCheck: DefaultHealthCheckConfig(),
		logger:      slog.Default(),
	}, nil
}

const expectedHostPortParts = 2

// Register registers a service instance with Consul including gRPC health checks.
func (r *ConsulRegistry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != expectedHostPortParts {
		return fmt.Errorf("invalid host:port format : %s", hostPort)
	}

	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	registration := &consulapi.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: host,
		Port:    port,
		Tags:    []string{"grpc"},
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", host, port),
			Interval:                       r.healthCheck.Interval,
			Timeout:                        r.healthCheck.Timeout,
			DeregisterCriticalServiceAfter: r.healthCheck.DeregisterAfter,
		},
	}

	if regErr := r.client.Agent().ServiceRegister(registration); regErr != nil {
		return regErr
	}

	r.logger.InfoContext(ctx, "🎉 Registered service", "serviceName", serviceName, "instanceID", instanceID)
	return nil
}

// Deregister removes a service instance from Consul.
func (r *ConsulRegistry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	if err := r.client.Agent().ServiceDeregister(instanceID); err != nil {
		return err
	}

	r.logger.InfoContext(ctx, "🔌 Deregistered service", "serviceName", serviceName, "instanceID", instanceID)
	return nil
}

// Connect establishes a gRPC connection to a service via Consul with load balancing.
func (r *ConsulRegistry) Connect(ctx context.Context, hostPort string, serviceName string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s/%s?tag=grpc&healthy=true", hostPort, serviceName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		return nil, err
	}

	r.logger.InfoContext(ctx, "🔗 Connected to service via Consul", "serviceName", serviceName)
	return conn, nil
}

// Close gracefully shuts down the registry (for interface consistency).
func (r *ConsulRegistry) Close() error {
	// Consul client doesn't require explicit cleanup
	return nil
}
