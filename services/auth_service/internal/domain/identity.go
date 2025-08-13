package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Identity represents an authentication provider connection for a user.
// It stores the mapping between a user and their identity from both external providers
// (like Google, Facebook, and other OAuth providers) and local email authentication.
type Identity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id"`
	ProviderID  string             `bson:"provider_id"`
	Provider    string             `bson:"provider"`
	Email       string             `bson:"email"`
	LastLoginAt time.Time          `bson:"last_login_at"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// IdentityRepository defines the interface for user identity data persistence operations.
type IdentityRepository interface {
	CreateIdentity(ctx context.Context, identity *Identity) (*Identity, error)
	GetIdentitiesByUserID(ctx context.Context, userID string) ([]Identity, error)
	GetIdentityByProvider(ctx context.Context, providerID string, provider string) (*Identity, error)
	UpdateLastLogin(ctx context.Context, userID string) error
}
