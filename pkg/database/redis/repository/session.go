package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ayehia0/org/pkg/database/mongodb/models"
	"github.com/redis/go-redis/v9"
)

// the session repository for managing the sessions using redis
type SessionRepository interface {
	CreateSession(cts context.Context, session *models.Session) error
	GetSessionByID(ctx context.Context, id string) (*models.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

type sessionRepository struct {
	conn *redis.Conn
}

func NewSessionRepository(conn *redis.Conn) SessionRepository {
	return &sessionRepository{conn: conn}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *models.Session) error {
	// add a session to the redis database
	// convert the session to a json
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// add the session to the redis database
	err = r.conn.Set(ctx, session.ID, sessionJSON, session.RefreshTokenExpires.Sub(time.Now())).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id string) (*models.Session, error) {
	// get the session from the redis database
	sessionJSON, err := r.conn.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}

	// convert the session to a struct
	var session models.Session
	err = json.Unmarshal([]byte(sessionJSON), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) DeleteSession(ctx context.Context, id string) error {
	// delete the session from the redis database
	err := r.conn.Del(ctx, id).Err()
	if err != nil {
		return err
	}

	return nil
}
