package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// here we setup the connection to the redis database and return the connection

type RedisConn struct {
	Conn *redis.Conn
}

func NewRedisConn(host string, port int, password string, database int) (*RedisConn, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       database,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	conn := client.Conn()
	return &RedisConn{Conn: conn}, nil
}
