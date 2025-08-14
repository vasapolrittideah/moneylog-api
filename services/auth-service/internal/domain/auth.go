package domain

import (
	"context"

	authtypes "github.com/vasapolrittideah/moneylog-api/services/auth-service/pkg/types"
)

// AuthUsecase defines the interface for authentication business logics.
type AuthUsecase interface {
	Login(ctx context.Context, params LoginParams) (*authtypes.Tokens, error)
	SignUp(ctx context.Context, params SignUpParams) (*authtypes.Tokens, error)
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
