package controllers

import (
	"net/http"

	api "github.com/ayehia0/org/pkg/api/middleware"
	"github.com/ayehia0/org/pkg/database/mongodb"
	"github.com/ayehia0/org/pkg/database/mongodb/models"
	"github.com/ayehia0/org/pkg/token"
	"github.com/ayehia0/org/pkg/utils"
	"github.com/gin-gonic/gin"
)

type OrgController interface {
	CreateOrganizationController(ctx *gin.Context)       // create a new organization
	DeleteOrganizationController(ctx *gin.Context)       // delete an organization
	UpdateOrganizationController(ctx *gin.Context)       // update an organization
	GetAllOrganizationsController(ctx *gin.Context)      // get all organizations
	GetOrganizationByIDController(ctx *gin.Context)      // get an organization by id
	InviteUserToOrganizationController(ctx *gin.Context) // invite a user to an organization
}

type orgController struct {
	MongoDBConn  *mongodb.MongoDBConn
	Store        *mongodb.Store
	TokenCreator token.TokenCreator
	AppConfig    *utils.AppConfig
}

func NewOrgController(conn *mongodb.MongoDBConn, token token.TokenCreator, appConfig *utils.AppConfig, store *mongodb.Store) OrgController {
	return &orgController{
		MongoDBConn:  conn,
		TokenCreator: token,
		AppConfig:    appConfig,
		Store:        store,
	}
}

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"description" binding:"required"`
}

func (o *orgController) CreateOrganizationController(ctx *gin.Context) {
	var req CreateOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResp(err))
		return
	}

	payload := ctx.MustGet(api.AuthPayloadKey).(*token.Payload)
	id, err := o.Store.OrganizationRepository.Create(ctx, &models.Organization{
		Name:    req.Name,
		Desc:    req.Desc,
		Members: []models.Member{},
		Creator: payload.UserId,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":         "Organization has been created successfully",
		"organization_id": id,
	})
}

func (o *orgController) DeleteOrganizationController(ctx *gin.Context) {
	id := ctx.Param("id")
	org, err := o.Store.OrganizationRepository.FindByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResp(err))
		return
	}

	payload := ctx.MustGet(api.AuthPayloadKey).(*token.Payload)

	// check if the authenticated user is the creator of the organization
	if org.Creator != payload.UserId {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResp(err))
		return
	}

	err = o.Store.OrganizationRepository.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Organization has been deleted successfully",
	})
}

type UpdateOrganizationRequest struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func (o *orgController) UpdateOrganizationController(ctx *gin.Context) {
	// update an organization
	var req UpdateOrganizationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResp(err))
		return
	}

	// get the organization by id
	id := ctx.Param("id")

	org, err := o.Store.OrganizationRepository.Update(ctx, &models.Organization{
		ID:   id,
		Name: req.Name,
		Desc: req.Desc,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"organization_id": org.ID,
		"name":            org.Name,
		"description":     org.Desc,
	})
}

// TODO: don't return the the creator nor id
func (o *orgController) GetAllOrganizationsController(ctx *gin.Context) {
	orgs, err := o.Store.OrganizationRepository.FindAll(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, orgs)
}

func (o *orgController) GetOrganizationByIDController(ctx *gin.Context) {
	id := ctx.Param("id")

	org, err := o.Store.OrganizationRepository.FindByID(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"organization_id": org.ID,
		"name":            org.Name,
		"description":     org.Desc,
		"members":         org.Members,
	})
}

type InviteUserToOrganizationRequest struct {
	Email string `json:"user_email" binding:"required,email"`
}

func (o *orgController) InviteUserToOrganizationController(ctx *gin.Context) {
	var req InviteUserToOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResp(err))
		return
	}

	id := ctx.Param("id")

	org, err := o.Store.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResp(err))
		return
	}

	// check if the user is already a member of the organization
	isMember, err := o.Store.OrganizationRepository.IsUserInOrganization(ctx, id, req.Email)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(err))
		return
	}

	if isMember {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User is already a member of the organization",
		})
		return
	}

	member := models.Member{
		ID:          org.ID,
		Name:        org.Name,
		Email:       org.Email,
		AccessLevel: "member", // you can change the access level
	}

	err = o.Store.OrganizationRepository.InviteUserToOrganization(ctx, id, member)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User has been invited to the organization successfully",
	})
}
