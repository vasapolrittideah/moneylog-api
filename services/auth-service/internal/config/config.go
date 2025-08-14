package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
)

type AuthServiceConfig struct {
	Environment string `env:"ENVIRONMENT"`
	Name        string `env:"AUTH_SERVICE_NAME"`
	Port        string `env:"AUTH_SERVICE_PORT"`
	Token       TokenConfig
}

type TokenConfig struct {
	AccessTokenSecret     string        `env:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret    string        `env:"REFRESH_TOKEN_SECRET"`
	AccessTokenExpiresIn  time.Duration `env:"ACCESS_TOKEN_EXPIRES_IN"`
	RefreshTokenExpiresIn time.Duration `env:"REFRESH_TOKEN_EXPIRES_IN"`
	Issuer                string        `env:"TOKEN_ISSUER"`
}

func NewAuthServiceConfig(logger *zerolog.Logger) *AuthServiceConfig {
	cfg, err := env.ParseAs[AuthServiceConfig]()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse env")
	}

	return &cfg
}
