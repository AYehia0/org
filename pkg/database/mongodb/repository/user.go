package repository

import (
	"context"
	"errors"

	"github.com/ayehia0/org/pkg/database/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// the repository package contains the database operations for the user model
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error                 // Create a new user
	FindByEmail(ctx context.Context, email string) (*models.User, error) // Find a user by email
	FindByID(ctx context.Context, id string) (*models.User, error)       // Find a user by id
}

// create a new user repository
func NewUserRepository(col *mongo.Collection) UserRepository {
	return &userRepository{col: col}
}

// the user repository struct
type userRepository struct {
	col *mongo.Collection
}

// the function to create a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	// make sure that the email is unique
	_, err := r.FindByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("email already exists")
	}
	_, err = r.col.InsertOne(ctx, user)
	return err
}

// the function to find a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, err
}

// the function to find a user by id
func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, err
}
