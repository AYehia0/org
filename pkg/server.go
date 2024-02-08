package pkg

import (
	"fmt"

	"github.com/ayehia0/org/pkg/api/handlers"
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

	// setup user routes
	userHandler := handlers.NewUserHandler(s.MongoDBConn, tokenCreator, s.AppConfig, s.Store)

	routes.SetupUserRoutes(s.Router.Group("/"), userHandler)

	// the api
	return nil
}

func (s *Server) Run() error {
	// run the server
	fmt.Printf("Server running on port %d\n", s.AppConfig.Port)
	return s.Router.Run(fmt.Sprintf("0.0.0.0:%d", s.AppConfig.Port))
}
