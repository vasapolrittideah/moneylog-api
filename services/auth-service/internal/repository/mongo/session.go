package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/vasapolrittideah/moneylog-api/services/auth-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const sessionCollection = "sessions"

type sessionMongoRepository struct {
	db *mongo.Database
}

func NewSessionRepository(_ context.Context, _ *zerolog.Logger, db *mongo.Database) domain.SessionRepository {
	return &sessionMongoRepository{
		db: db,
	}
}

func (r *sessionMongoRepository) CreateSession(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	now := time.Now()
	session.CreatedAt = now
	session.UpdatedAt = now

	result, err := r.db.Collection(sessionCollection).InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert inserted ID to ObjectID")
	}
	session.ID = objectID

	return session, nil
}

func (r *sessionMongoRepository) GetSessionByUserID(ctx context.Context, userID string) (*domain.Session, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	result := r.db.Collection(sessionCollection).FindOne(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var session domain.Session
	if err := result.Decode(&session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionMongoRepository) UpdateTokens(
	ctx context.Context,
	id string,
	params domain.UpdateTokensParams,
) (*domain.Session, error) {
	result := r.db.Collection(sessionCollection).FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": params},
	)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var session domain.Session
	if err := result.Decode(&session); err != nil {
		return nil, err
	}

	return &session, nil
}
