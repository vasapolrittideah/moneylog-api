package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/vasapolrittideah/moneylog-api/services/auth_service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const userCollection = "users"

type userMongoRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &userMongoRepository{
		db: db,
	}
}

func (r *userMongoRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(userCollection)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *userMongoRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := r.db.Collection(userCollection).InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert inserted ID to ObjectID")
	}
	user.ID = objectID

	return user, nil
}

func (r *userMongoRepository) GetUser(ctx context.Context, id string) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := r.db.Collection(userCollection).FindOne(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var user domain.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userMongoRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	result := r.db.Collection(userCollection).FindOne(ctx, bson.M{"email": email})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var user domain.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userMongoRepository) UpdateUser(
	ctx context.Context,
	id string,
	params domain.UpdateUserParams,
) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Build update query
	updateMap := bson.M{}
	if params.Email != nil {
		updateMap["email"] = params.Email
	}
	if params.FullName != nil {
		updateMap["full_name"] = params.FullName
	}
	if params.PasswordHash != nil {
		updateMap["password_hash"] = params.PasswordHash
	}

	if len(updateMap) == 0 {
		return nil, errors.New("no user fields to update")
	}

	updateMap["updated_at"] = time.Now()

	result := r.db.Collection(userCollection).FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateMap},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var user domain.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userMongoRepository) DeleteUser(ctx context.Context, id string) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := r.db.Collection(userCollection).FindOneAndDelete(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var user domain.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userMongoRepository) ListUsers(ctx context.Context, params domain.FilterUserParams) ([]*domain.User, error) {
	findOptions := options.Find()

	limit := params.Limit
	if limit == 0 {
		limit = 10
	}
	findOptions.SetLimit(int64(limit))

	if params.Offset > 0 {
		findOptions.SetSkip(int64(params.Offset))
	}

	sortBy := "created_at"
	if params.SortBy != nil {
		sortBy = *params.SortBy
	}

	sortOrder := -1
	if !params.SortDesc {
		sortOrder = 1
	}
	findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	// Build filter query
	filter := bson.M{}
	if params.Email != nil {
		filter["email"] = *params.Email
	}
	if params.Verified != nil {
		filter["verified"] = *params.Verified
	}

	cursor, err := r.db.Collection(userCollection).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
