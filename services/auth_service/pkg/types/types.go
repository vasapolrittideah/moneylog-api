package authtypes

import "github.com/golang-jwt/jwt/v5"

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type JWTClaims struct {
	jwt.RegisteredClaims

	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
}
