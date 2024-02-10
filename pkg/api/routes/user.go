package routes

import (
	"github.com/ayehia0/org/pkg/api/handlers"
	"github.com/gin-gonic/gin"
)

// here we define all the routes for user business logic
func SetupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler) {
	router.POST("/signup", userHandler.SignupHandler)
	router.POST("/login", userHandler.LoginHandler)
	router.POST("/refresh-token", userHandler.RefreshTokenHandler)
	router.POST("/revoke-refresh-token", userHandler.RevokeRefreshTokenHandler)
}
