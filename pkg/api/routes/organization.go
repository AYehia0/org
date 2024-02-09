package routes

import (
	"github.com/ayehia0/org/pkg/api/handlers"
	"github.com/gin-gonic/gin"
)

// here we define all the routes for user business logic
func SetupOrgRoutes(router *gin.RouterGroup, orgHandler *handlers.OrgHandler) {
	router.POST("/", orgHandler.CreateOrganizationHandler)
	router.POST("/:id/invite", orgHandler.InviteUserToOrganizationHandler)
	router.DELETE("/:id", orgHandler.DeleteOrganizationHandler)
	router.PUT("/:id", orgHandler.UpdateOrganizationHandler)
	router.GET("/", orgHandler.GetAllOrganizationsHandler)
	router.GET("/:id", orgHandler.GetOrganizationByIDHandler)
}
