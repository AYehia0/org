package handlers

import (
	"github.com/ayehia0/org/pkg/controllers"
	"github.com/ayehia0/org/pkg/database/mongodb"
	"github.com/ayehia0/org/pkg/token"
	"github.com/ayehia0/org/pkg/utils"
	"github.com/gin-gonic/gin"
)

// here we define all the handlers for user business logic like signup, login, etc.
type UserHandler interface {
	SignupHandler(ctx *gin.Context)
	LoginHandler(ctx *gin.Context)
	RefreshTokenHandler(ctx *gin.Context)
}

type userHandler struct {
	// the handler should have a reference to the controller
	userController controllers.UserController
}

func NewUserHandler(conn *mongodb.MongoDBConn, token token.TokenCreator, appConfig *utils.AppConfig, store *mongodb.Store) UserHandler {
	return &userHandler{
		userController: controllers.NewUserController(conn, token, appConfig, store),
	}
}

// the handlers should return a gin.HandlerFunc
func (u *userHandler) SignupHandler(ctx *gin.Context) {
	u.userController.SignupController(ctx)
}

func (u *userHandler) LoginHandler(ctx *gin.Context) {
	u.userController.LoginController(ctx)
}

func (u *userHandler) RefreshTokenHandler(ctx *gin.Context) {
	u.userController.RefreshTokenController(ctx)
}
