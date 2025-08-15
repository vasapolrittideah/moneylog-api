package authclient

import (
	"fmt"

	"github.com/vasapolrittideah/moneylog-api/services/auth-service/internal/config"
	"github.com/vasapolrittideah/moneylog-api/shared/discovery"
	authpbv1 "github.com/vasapolrittideah/moneylog-api/shared/protos/auth/v1"
	"google.golang.org/grpc"
)

type AuthServiceClient struct {
	Client authpbv1.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthServiceClient(
	consulRegistry *discovery.ConsulRegistry,
	authServiceCfg *config.AuthServiceConfig,
) (*AuthServiceClient, error) {
	authServiceAddr := fmt.Sprintf("%s:%s", authServiceCfg.Host, authServiceCfg.Port)
	conn, err := consulRegistry.Connect(authServiceAddr, authServiceCfg.Name)
	if err != nil {
		return nil, err
	}

	client := authpbv1.NewAuthServiceClient(conn)

	return &AuthServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *AuthServiceClient) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
