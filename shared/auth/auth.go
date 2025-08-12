package auth

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateToken(claims jwt.Claims, secret string) (string, error)
	ValidateToken(token string, secret string) (*jwt.Token, error)
}
