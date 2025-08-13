package domain

import (
	"context"

	"github.com/vasapolrittideah/moneylog-api/services/auth_service/pkg/types"
)

// AuthUsecase defines the interface for authentication business logics.
type AuthUsecase interface {
	Login(ctx context.Context, params LoginParams) (*types.Tokens, error)
	SignUp(ctx context.Context, params SignUpParams) (*types.Tokens, error)
}

// LoginParams contains the parameters for user login.
type LoginParams struct {
	Email    string
	Password string
}

// SignUpParams contains the parameters for user sign up.
type SignUpParams struct {
	Email    string
	Password string
	FullName string
}
