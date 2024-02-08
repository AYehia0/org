package controllers

import (
	"errors"
	"net/http"

	"github.com/ayehia0/org/pkg/database/mongodb"
	"github.com/ayehia0/org/pkg/database/mongodb/models"
	"github.com/ayehia0/org/pkg/token"
	"github.com/ayehia0/org/pkg/utils"
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

	// contains all the repositories
	Store *mongodb.Store

	// the token creator used to create and verify the tokens
	TokenCreator token.TokenCreator

	// TODO: Change this : there is a better way
	AppConfig *utils.AppConfig
}

func NewUserController(conn *mongodb.MongoDBConn, token token.TokenCreator, appConfig *utils.AppConfig, store *mongodb.Store) UserController {
	return &userController{
		MongoDBConn:  conn,
		TokenCreator: token,
		AppConfig:    appConfig,
		Store:        store,
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
		ctx.JSON(http.StatusBadRequest, utils.ErrorResp(err))
		return
	}
	// hashing the password
	password, err := utils.GenerateHash(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to hash the password")))
		return
	}

	// save the user to the database
	err = u.Store.UserRepository.Create(ctx, &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(err))
		return
	}

	// for testing return the request
	ctx.JSON(http.StatusOK, gin.H{"message": "User has been created successfully"})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (u *userController) LoginController(ctx *gin.Context) {
	var req LoginRequest

	// bind the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResp(err))
		return
	}

	// get the user from the Database
	user, err := u.Store.UserRepository.FindByEmail(ctx, req.Email)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to get the user")))
		return
	}

	// compare the password
	if err := utils.ComparePasswords(req.Password, user.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResp(errors.New("invalid credentials")))
		return
	}

	// create a access token
	userAccessToken, payloadAccess, err := u.TokenCreator.Create(user.Email, u.AppConfig.TokenAccessExpiration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to create a token")))
		return
	}

	// create a refresh token
	userRefreshToken, payloadrefresh, err := u.TokenCreator.Create(user.Email, u.AppConfig.TokenRefreshExpiration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to create a token")))
		return
	}

	// create the session to the database
	err = u.Store.SessionRepository.Create(ctx, &models.Session{
		UserID:              user.ID,
		AccessToken:         userAccessToken,
		RefreshToken:        userRefreshToken,
		AccessTokenExpires:  payloadAccess.ExpiredAt,
		RefreshTokenExpires: payloadrefresh.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to create a session")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "User has been logged in successfully",
		"access_token":  userAccessToken,
		"refresh_token": userRefreshToken,
	})
}
