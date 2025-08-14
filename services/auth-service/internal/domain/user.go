package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user account in the authentication system.
type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	FullName         string             `bson:"full_name"`
	Email            string             `bson:"email"`
	PasswordHash     string             `bson:"password_hash"`
	Verified         bool               `bson:"verified"`
	VerificationCode string             `bson:"verification_code"`
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
}

// UserRepository defines the interface for user data persistence operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id string, params UpdateUserParams) (*User, error)
	DeleteUser(ctx context.Context, id string) (*User, error)
	ListUsers(ctx context.Context, params FilterUserParams) ([]*User, error)
}

// UpdateUserParams contains the optional parameters for updating a user.
// Only non-nil fields will be updated.
type UpdateUserParams struct {
	Email        *string
	FullName     *string
	PasswordHash *string
}

// FilterUserParams contains the parameters for filtering and paginating user queries.
type FilterUserParams struct {
	Email    *string
	Verified *bool
	Limit    uint64
	Offset   uint64
	SortBy   *string
	SortDesc bool
}
