package pkg

import (
	"fmt"

	types "github.com/ayehia0/org/pkg/api"
	"github.com/ayehia0/org/pkg/api/handlers"
	api "github.com/ayehia0/org/pkg/api/middleware"
	"github.com/ayehia0/org/pkg/api/routes"
	"github.com/ayehia0/org/pkg/database/mongodb"
	"github.com/ayehia0/org/pkg/database/redis"
	"github.com/ayehia0/org/pkg/token"
	"github.com/ayehia0/org/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	MongoDBConn *mongodb.MongoDBConn
	RedisConn   *redis.RedisConn
	DBConfig    *utils.DatabaseConfig
	AppConfig   *utils.AppConfig
	RedisConfig *utils.RedisConfig
	Router      *gin.Engine
	DBStore     *mongodb.DBStore
	RedisStore  *redis.RedisStore
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	// the configs
	dbConfig, redisConfig, appConfig, err := utils.ConfigStore("./config", "database-config", "redis-config", "app-config")

	if err != nil {
		return err
	}

	s.DBConfig = &dbConfig
	s.AppConfig = &appConfig
	s.RedisConfig = &redisConfig

	// connect to the database
	uri := fmt.Sprintf("mongodb://%s:%d", s.DBConfig.Host, s.DBConfig.Port)
	dbConn, err := mongodb.NewMongoDBConn(uri, s.DBConfig.Database, s.DBConfig.Username, s.DBConfig.Password)
	if err != nil {
		return err
	}

	// connect to the redis
	redisDb := 0 // the default db
	redisConn, err := redis.NewRedisConn(s.RedisConfig.Host, s.RedisConfig.Port, s.RedisConfig.Password, redisDb)

	s.MongoDBConn = dbConn
	s.RedisConn = redisConn

	// setup the engine
	s.Router = gin.Default()

	// defining the repositories
	s.DBStore = mongodb.NewStore(s.MongoDBConn)
	s.RedisStore = redis.NewStore(s.RedisConn)

	// create a token creator
	tokenCreator, err := token.NewPasteoToken(s.AppConfig.JwtSecret)

	if err != nil {
		return err
	}

	appC := &types.AppC{
		MongoDBConn:  s.MongoDBConn,
		DBStore:      s.DBStore,
		RDBStore:     s.RedisStore,
		TokenCreator: tokenCreator,
		AppConfig:    s.AppConfig,
	}

	orgHandler := handlers.NewOrgHandler(appC)
	userHandler := handlers.NewUserHandler(appC)

	routes.SetupUserRoutes(s.Router.Group("/"), userHandler)

	// use authMiddleware to protect the routes
	authRquired := s.Router.Group("/organizations")
	authRquired.Use(api.AuthMiddleware(tokenCreator))

	routes.SetupOrgRoutes(authRquired, orgHandler)

	// the api
	return nil
}

func (s *Server) Run() error {
	// run the server

	// check if the server is running on production or development
	if s.AppConfig.Env == "production" {
		// run a ssl server using the certs issued by letsencrypt which is found on : /etc/letsencrypt/live/<domain-name>/{fullchain.pem, privkey.pem}
		// disable debug mode
		gin.SetMode(gin.ReleaseMode)
		return s.Router.RunTLS(fmt.Sprintf("0.0.0.0:%d", s.AppConfig.Port),
			"./fullchain.pem",
			"./privkey.pem",
		)
	} else {
		return s.Router.Run(fmt.Sprintf("0.0.0.0:%d", s.AppConfig.Port))
	}
}
