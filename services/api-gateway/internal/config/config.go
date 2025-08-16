package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
)

type APIGatewayConfig struct {
	Environment string `env:"ENVIRONMENT"`
	Addr        string `env:"API_GATEWAY_ADDR"`
	AuthService AuthServiceConfig
}

type AuthServiceConfig struct {
	Name string `env:"AUTH_SERVICE_NAME"`
}

func NewAPIGatewayConfig(logger *zerolog.Logger) *APIGatewayConfig {
	cfg, err := env.ParseAs[APIGatewayConfig]()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse env")
	}

	return &cfg
}
