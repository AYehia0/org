package pkg

import (
	"fmt"

	types "github.com/ayehia0/org/pkg/api"
	"github.com/ayehia0/org/pkg/api/handlers"
	api "github.com/ayehia0/org/pkg/api/middleware"
	"github.com/ayehia0/org/pkg/api/routes"
	"github.com/ayehia0/org/pkg/database/mongodb"
	"github.com/ayehia0/org/pkg/token"
	"github.com/ayehia0/org/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	MongoDBConn *mongodb.MongoDBConn
	DBConfig    *utils.DatabaseConfig
	AppConfig   *utils.AppConfig
	Router      *gin.Engine
	Store       *mongodb.Store
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	// the configs
	dbConfig, appConfig, err := utils.ConfigStore("./config", "database-config", "app-config")

	if err != nil {
		return err
	}

	s.DBConfig = &dbConfig
	s.AppConfig = &appConfig

	// connect to the database
	uri := fmt.Sprintf("mongodb://%s:%d", s.DBConfig.Host, s.DBConfig.Port)
	conn, err := mongodb.NewMongoDBConn(uri, s.DBConfig.Database, s.DBConfig.Username, s.DBConfig.Password)
	if err != nil {
		return err
	}

	s.MongoDBConn = conn

	// setup the engine
	s.Router = gin.Default()

	// defining the repositories
	s.Store = mongodb.NewStore(s.MongoDBConn)

	// create a token creator
	tokenCreator, err := token.NewPasteoToken(s.AppConfig.JwtSecret)

	if err != nil {
		return err
	}

	appC := &types.AppC{
		MongoDBConn:  s.MongoDBConn,
		Store:        s.Store,
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
		return s.Router.RunTLS(fmt.Sprintf("0.0.0.0:%d", s.AppConfig.Port),
			"./fullchain.pem",
			"./privkey.pem",
		)
	} else {
		fmt.Printf("Server running on port %d\n", s.AppConfig.Port)
		return s.Router.Run(fmt.Sprintf("0.0.0.0:%d", s.AppConfig.Port))
	}
}
