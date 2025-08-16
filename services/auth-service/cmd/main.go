package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/vasapolrittideah/moneylog-api/services/auth-service/internal/config"
	grpchandler "github.com/vasapolrittideah/moneylog-api/services/auth-service/internal/delivery/grpc"
	mongodb "github.com/vasapolrittideah/moneylog-api/services/auth-service/internal/repository/mongo"
	"github.com/vasapolrittideah/moneylog-api/services/auth-service/internal/usecase"
	"github.com/vasapolrittideah/moneylog-api/shared/auth"
	"github.com/vasapolrittideah/moneylog-api/shared/database"
	"github.com/vasapolrittideah/moneylog-api/shared/discovery"
	"github.com/vasapolrittideah/moneylog-api/shared/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logger.Get()

	authServiceCfg := config.NewAuthServiceConfig(logger)
	mongoCfg := database.NewMongoConfig(logger)

	mongoDB := database.NewMongoDB(mongoCfg, logger)
	if err := mongoDB.Connect(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	defer func() {
		if err := mongoDB.Disconnect(ctx); err != nil {
			logger.Error().Err(err).Msg("Failed to disconnect from MongoDB")
		}
	}()

	consulCfg := discovery.NewConsulRegistryConfig(logger)
	consulRegistry, err := discovery.NewConsulRegistry(consulCfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create Consul registry")
	}

	serviceID := authServiceCfg.Name + "-1"
	if err := consulRegistry.Register(serviceID, authServiceCfg.Name, authServiceCfg.RegisterAddr); err != nil {
		logger.Fatal().Err(err).Msg("Failed to register service in Consul")
	}
	defer func() {
		if err := consulRegistry.Deregister(serviceID, authServiceCfg.Name); err != nil {
			logger.Error().Err(err).Msg("Failed to deregister service in Consul")
		}
	}()

	jwtAuthenticator := auth.NewJWTAuthenticator(
		authServiceCfg.Token.Issuer,
		authServiceCfg.Token.Issuer,
	)

	identityRepo := mongodb.NewIdentityRepository(ctx, logger, mongoDB.GetDatabase())
	sessionRepo := mongodb.NewSessionRepository(ctx, logger, mongoDB.GetDatabase())
	userRepo := mongodb.NewUserRepository(ctx, logger, mongoDB.GetDatabase())

	authUsecase := usecase.NewAuthUsecase(identityRepo, sessionRepo, userRepo, jwtAuthenticator, authServiceCfg)

	lc := net.ListenConfig{}
	lis, err := lc.Listen(ctx, "tcp", authServiceCfg.Addr)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to listen on gRPC address")
	}

	grpcServer := grpc.NewServer()
	grpchandler.NewAuthGRPCHandler(grpcServer, authUsecase)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan
		cancel()
	}()

	logger.Info().Str("address", authServiceCfg.Addr).Msg("Starting gRPC server for auth service")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start gRPC server")
			cancel()
		}
	}()

	<-ctx.Done()
	logger.Info().Msg("Shutting down gRPC server")
	grpcServer.GracefulStop()
}
