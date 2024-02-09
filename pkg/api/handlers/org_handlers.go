package handlers

import (
	"github.com/ayehia0/org/pkg/controllers"
	"github.com/ayehia0/org/pkg/database/mongodb"
	"github.com/ayehia0/org/pkg/token"
	"github.com/ayehia0/org/pkg/utils"
	"github.com/gin-gonic/gin"
)

type OrganizationHandler interface {
	CreateOrganizationHandler(ctx *gin.Context)       // create a new organization
	DeleteOrganizationHandler(ctx *gin.Context)       // delete an organization
	UpdateOrganizationHandler(ctx *gin.Context)       // update an organization
	GetAllOrganizationsHandler(ctx *gin.Context)      // get all organizations
	GetOrganizationByIDHandler(ctx *gin.Context)      // get an organization by id
	InviteUserToOrganizationHandler(ctx *gin.Context) // invite a user to an organization
}

type orgHandler struct {
	// the handler should have a reference to the controller
	OrgController controllers.OrgController
}

func NewOrgHandler(conn *mongodb.MongoDBConn, token token.TokenCreator, appConfig *utils.AppConfig, store *mongodb.Store) OrganizationHandler {
	return &orgHandler{
		OrgController: controllers.NewOrgController(conn, token, appConfig, store),
	}
}

func (o *orgHandler) CreateOrganizationHandler(ctx *gin.Context) {
	o.OrgController.CreateOrganizationController(ctx)
}

func (o *orgHandler) DeleteOrganizationHandler(ctx *gin.Context) {
	o.OrgController.DeleteOrganizationController(ctx)
}

func (o *orgHandler) UpdateOrganizationHandler(ctx *gin.Context) {
	o.OrgController.UpdateOrganizationController(ctx)
}

func (o *orgHandler) GetAllOrganizationsHandler(ctx *gin.Context) {
	o.OrgController.GetAllOrganizationsController(ctx)
}

func (o *orgHandler) GetOrganizationByIDHandler(ctx *gin.Context) {
	o.OrgController.GetOrganizationByIDController(ctx)
}

func (o *orgHandler) InviteUserToOrganizationHandler(ctx *gin.Context) {
	o.OrgController.InviteUserToOrganizationController(ctx)
}
