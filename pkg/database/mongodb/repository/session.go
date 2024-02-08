package repository

import (
	"context"
	"errors"

	"github.com/ayehia0/org/pkg/database/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// the repository package contains the database operations for the session model
type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error                // Create a new Session
	FindByID(ctx context.Context, id string) (*models.Session, error)         // Find a Session by id
	FindByUserID(ctx context.Context, userID string) (*models.Session, error) // Find a Session by user id
	Delete(ctx context.Context, id string) error                              // Delete a Session
}

// the session repository struct
type sessionRepository struct {
	// the database connection
	col *mongo.Collection
}

// create a new session repository
func NewSessionRepository(col *mongo.Collection) SessionRepository {
	return &sessionRepository{col: col}
}

// the function to create a new session
func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	_, err := r.col.InsertOne(ctx, session)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("session already exists")
		}
	}
	return err
}

// the function to find a session by id
func (r *sessionRepository) FindByID(ctx context.Context, id string) (*models.Session, error) {
	var session models.Session
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("session not found")
		}
	}
	return &session, err
}

// the function to find a session by user id
func (r *sessionRepository) FindByUserID(ctx context.Context, userID string) (*models.Session, error) {
	var session models.Session
	err := r.col.FindOne(ctx, bson.M{"user_id": userID}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("session not found")
		}
	}
	return &session, err
}

// the function to delete a session
func (r *sessionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("session not found")
		}
	}
	return err
}
