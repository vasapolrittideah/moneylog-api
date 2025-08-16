package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/moneylog-api/services/api-gateway/internal/config"
	httphandler "github.com/vasapolrittideah/moneylog-api/services/api-gateway/internal/delivery/http"
	authclient "github.com/vasapolrittideah/moneylog-api/services/auth-service/pkg/client"
	"github.com/vasapolrittideah/moneylog-api/shared/discovery"
	"github.com/vasapolrittideah/moneylog-api/shared/logger"
)

const appShutdownTimout = 10 * time.Second

func main() {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	logger := logger.Get()

	consulCfg := discovery.NewConsulRegistryConfig(logger)
	consulRegistry, err := discovery.NewConsulRegistry(consulCfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create Consul registry")
	}

	apiGatewayCfg := config.NewAPIGatewayConfig(logger)
	authServiceClient, err := authclient.NewAuthServiceClient(
		consulRegistry,
		apiGatewayCfg.AuthService.Name,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create auth client")
	}

	app := fiber.New()

	authHandler := httphandler.NewAuthHTTPHandler(authServiceClient, app, logger)
	authHandler.RegisterRoutes()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info().Str("address", apiGatewayCfg.Addr).Msg("Starting API Gateway")
		serverErrors <- app.Listen(apiGatewayCfg.Addr)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error().Err(err).Msg("Failed to start server")

	case sig := <-shutdown:
		logger.Info().Interface("signal", sig).Msg("Received shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), appShutdownTimout)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			logger.Error().Err(err).Msg("Failed to shutdown server")
		}
	}
}
