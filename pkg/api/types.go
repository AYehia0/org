package types

import (
	"github.com/ayehia0/org/pkg/database/mongodb"
	"github.com/ayehia0/org/pkg/database/redis"
	"github.com/ayehia0/org/pkg/token"
	"github.com/ayehia0/org/pkg/utils"
)

type AppC struct {
	MongoDBConn  *mongodb.MongoDBConn
	DBStore      *mongodb.DBStore
	RDBStore     *redis.RedisStore
	TokenCreator token.TokenCreator
	AppConfig    *utils.AppConfig
}
