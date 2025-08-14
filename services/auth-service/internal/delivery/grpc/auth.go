package grpc

import (
	"context"
	"errors"

	"github.com/vasapolrittideah/moneylog-api/services/auth-service/internal/domain"
	"github.com/vasapolrittideah/moneylog-api/services/auth-service/internal/usecase"
	authpbv1 "github.com/vasapolrittideah/moneylog-api/shared/protos/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authGRPCHandler struct {
	authpbv1.UnimplementedAuthServiceServer

	authUsecase domain.AuthUsecase
}

func NewAuthGRPCHandler(server *grpc.Server, authUsecase domain.AuthUsecase) authpbv1.AuthServiceServer {
	handler := &authGRPCHandler{
		authUsecase: authUsecase,
	}
	authpbv1.RegisterAuthServiceServer(server, handler)

	return handler
}

func (h *authGRPCHandler) Login(ctx context.Context, req *authpbv1.LoginRequest) (*authpbv1.LoginResponse, error) {
	params := domain.LoginParams{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	tokens, err := h.authUsecase.Login(ctx, params)
	if err != nil {
		var code codes.Code
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			code = codes.Unauthenticated
		case errors.Is(err, usecase.ErrUserNotFound):
			code = codes.NotFound
		default:
			code = codes.Internal
		}

		return nil, status.Errorf(code, "failed to login: %v", err)
	}

	return &authpbv1.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *authGRPCHandler) SignUp(ctx context.Context, req *authpbv1.SignUpRequest) (*authpbv1.SignUpResponse, error) {
	params := domain.SignUpParams{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		FullName: req.GetFullName(),
	}

	tokens, err := h.authUsecase.SignUp(ctx, params)
	if err != nil {
		var code codes.Code
		switch {
		case errors.Is(err, usecase.ErrUserAlreadyExists):
			code = codes.AlreadyExists
		default:
			code = codes.Internal
		}

		return nil, status.Errorf(code, "failed to sign up: %v", err)
	}

	return &authpbv1.SignUpResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
