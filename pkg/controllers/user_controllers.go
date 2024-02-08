package controllers

import (
	"github.com/ayehia0/org/pkg/database/mongodb"
	"github.com/gin-gonic/gin"
)

// here we define all the controllers for user business logic like signup, login, etc.
type UserController interface {
	// the user can signup
	SignupController(ctx *gin.Context)

	// the user can login
	LoginController(ctx *gin.Context)
}

// the user controller should have access to the database
type userController struct {
	// the controller should have a reference to the database
	MongoDBConn *mongodb.MongoDBConn
}

func NewUserController(conn *mongodb.MongoDBConn) UserController {
	return &userController{
		MongoDBConn: conn,
	}
}

type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// the controllers should return a gin.HandlerFunc
func (u *userController) SignupController(ctx *gin.Context) {
	// call the controller to handle the request
	var req SignupRequest

	// bind the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// for testing return the request
	ctx.JSON(200, req)
}

func (u *userController) LoginController(ctx *gin.Context) {
}
