package handlers

import (
	types "github.com/ayehia0/org/pkg/api"
	"github.com/ayehia0/org/pkg/controllers"
	"github.com/gin-gonic/gin"
)

type OrgHandler struct {
	orgController controllers.OrgController
}

func NewOrgHandler(appC *types.AppC) *OrgHandler {
	orgController := controllers.NewOrgController(appC)
	return &OrgHandler{orgController: orgController}
}

func (o *OrgHandler) CreateOrganizationHandler(ctx *gin.Context) {
	o.orgController.CreateOrganizationController(ctx)
}

func (o *OrgHandler) DeleteOrganizationHandler(ctx *gin.Context) {
	o.orgController.DeleteOrganizationController(ctx)
}

func (o *OrgHandler) UpdateOrganizationHandler(ctx *gin.Context) {
	o.orgController.UpdateOrganizationController(ctx)
}

func (o *OrgHandler) GetAllOrganizationsHandler(ctx *gin.Context) {
	o.orgController.GetAllOrganizationsController(ctx)
}

func (o *OrgHandler) GetOrganizationByIDHandler(ctx *gin.Context) {
	o.orgController.GetOrganizationByIDController(ctx)
}

func (o *OrgHandler) InviteUserToOrganizationHandler(ctx *gin.Context) {
	o.orgController.InviteUserToOrganizationController(ctx)
}
