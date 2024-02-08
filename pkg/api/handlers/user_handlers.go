package handlers

import (
	"github.com/ayehia0/org/pkg/controllers"
	"github.com/ayehia0/org/pkg/database/mongodb"
	"github.com/gin-gonic/gin"
)

// here we define all the handlers for user business logic like signup, login, etc.
type UserHandler interface {
	// the user can signup
	SignupHandler(ctx *gin.Context)

	// the user can login
	LoginHandler(ctx *gin.Context)
}

type userHandler struct {
	// the handler should have a reference to the controller
	userController controllers.UserController
}

func NewUserHandler(conn *mongodb.MongoDBConn) UserHandler {
	return &userHandler{
		userController: controllers.NewUserController(conn),
	}
}

// the handlers should return a gin.HandlerFunc
func (u *userHandler) SignupHandler(ctx *gin.Context) {
	u.userController.SignupController(ctx)
}

func (u *userHandler) LoginHandler(ctx *gin.Context) {
	u.userController.LoginController(ctx)
}
