package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session represents an authenticated user session with access and refresh tokens.
// It tracks token expiration times and optional metadata like IP address and user agent
// for security and auditing purposes.
type Session struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"`
	UserID                string             `bson:"user_id"`
	AccessToken           string             `bson:"access_token"`
	RefreshToken          string             `bson:"refresh_token"`
	AccessTokenExpiresAt  time.Time          `bson:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time          `bson:"refresh_token_expires_at"`
	IPAddress             *string            `bson:"ip_address"`
	UserAgent             *string            `bson:"user_agent"`
	CreatedAt             time.Time          `bson:"created_at"`
	UpdatedAt             time.Time          `bson:"updated_at"`
}

// SessionRepository defines the interface for session data persistence operations.
type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session) (*Session, error)
	GetSessionByUserID(ctx context.Context, userID string) (*Session, error)
	UpdateTokens(ctx context.Context, userID string, params UpdateTokensParams) (*Session, error)
}

// UpdateTokensParams contains the parameters for updating session tokens.
type UpdateTokensParams struct {
	AccessToken           string    `bson:"access_token"`
	RefreshToken          string    `bson:"refresh_token"`
	AccessTokenExpiresAt  time.Time `bson:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `bson:"refresh_token_expires_at"`
}
