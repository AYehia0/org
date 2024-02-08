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
	ExtendsAccessToken(ctx context.Context, id string, expires int64) error   // Extends the access token expiration
	ExtendsRefreshToken(ctx context.Context, id string, expires int64) error  // Extends the refresh token expiration
}

// the session repository struct
type sessionRepository struct {
	// the database connection
	col *mongo.Collection
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
	err := r.col.FindOne(ctx, models.Session{ID: id}).Decode(&session)
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
	err := r.col.FindOne(ctx, models.Session{UserID: userID}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("session not found")
		}
	}
	return &session, err
}

// the function to delete a session
func (r *sessionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, models.Session{ID: id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("session not found")
		}
	}
	return err
}

// the function to extend the access token expiration
func (r *sessionRepository) ExtendsAccessToken(ctx context.Context, id string, expires int64) error {
	_, err := r.col.UpdateOne(ctx, models.Session{ID: id}, bson.M{"$set": bson.M{"access_token_expires": expires}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("session not found")
		}
	}
	return err
}

// the function to extend the refresh token expiration
func (r *sessionRepository) ExtendsRefreshToken(ctx context.Context, id string, expires int64) error {
	_, err := r.col.UpdateOne(ctx, models.Session{ID: id}, bson.M{"$set": bson.M{"refresh_token_expires": expires}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("session not found")
		}
	}
	return err
}
