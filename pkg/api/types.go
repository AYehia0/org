package types

import (
	"github.com/ayehia0/org/pkg/database/mongodb"
	"github.com/ayehia0/org/pkg/token"
	"github.com/ayehia0/org/pkg/utils"
)

type AppC struct {
	MongoDBConn  *mongodb.MongoDBConn
	Store        *mongodb.Store
	TokenCreator token.TokenCreator
	AppConfig    *utils.AppConfig
}
