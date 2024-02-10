package handlers

import (
	types "github.com/ayehia0/org/pkg/api"
	"github.com/ayehia0/org/pkg/controllers"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userController controllers.UserController
}

func NewUserHandler(appC *types.AppC) *UserHandler {
	userController := controllers.NewUserController(appC)
	return &UserHandler{userController: userController}
}

// the handlers should return a gin.HandlerFunc
func (u *UserHandler) SignupHandler(ctx *gin.Context) {
	u.userController.SignupController(ctx)
}

func (u *UserHandler) LoginHandler(ctx *gin.Context) {
	u.userController.LoginController(ctx)
}

func (u *UserHandler) RefreshTokenHandler(ctx *gin.Context) {
	u.userController.RefreshTokenController(ctx)
}

func (u *UserHandler) RevokeRefreshTokenHandler(ctx *gin.Context) {
	u.userController.RevokeRefreshTokenController(ctx)
}
