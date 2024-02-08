package routes

import (
	"github.com/ayehia0/org/pkg/api/handlers"
	"github.com/gin-gonic/gin"
)

// here we define all the routes for user business logic
func SetupUserRoutes(router *gin.RouterGroup, userHandler handlers.UserHandler) {

	// signup route
	router.POST("/signup", userHandler.SignupHandler)

	// login route
	router.POST("/login", userHandler.LoginHandler)
}
