package redis

import "github.com/ayehia0/org/pkg/database/redis/repository"

type RedisStore struct {
	SessionRepository repository.SessionRepository
}

func NewStore(conn *RedisConn) *RedisStore {
	session := repository.NewSessionRepository(conn.Conn)
	return &RedisStore{
		SessionRepository: session,
	}
}
