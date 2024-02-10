package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	types "github.com/ayehia0/org/pkg/api"
	"github.com/ayehia0/org/pkg/database/mongodb/models"
	"github.com/ayehia0/org/pkg/utils"
	"github.com/gin-gonic/gin"
)

// here we define all the controllers for user business logic like signup, login, etc.
type UserController interface {
	// the user can signup
	SignupController(ctx *gin.Context)

	// the user can login
	LoginController(ctx *gin.Context)

	// refresh token
	RefreshTokenController(ctx *gin.Context)
}

type appU struct {
	types.AppC
}

func NewUserController(appC *types.AppC) UserController {
	return &appU{AppC: *appC}
}

// the controllers should return a gin.HandlerFunc
func (au *appU) SignupController(ctx *gin.Context) {
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
	err = au.DBStore.UserRepository.Create(ctx, &models.User{
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

func (au *appU) LoginController(ctx *gin.Context) {
	var req LoginRequest

	// bind the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResp(err))
		return
	}

	// get the user from the Database
	user, err := au.DBStore.UserRepository.FindByEmail(ctx, req.Email)

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
	userAccessToken, payloadAccess, err := au.TokenCreator.Create(user.ID, au.AppConfig.TokenAccessExpiration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to create a token")))
		return
	}

	// create a refresh token
	userRefreshToken, payloadrefresh, err := au.TokenCreator.Create(user.ID, au.AppConfig.TokenRefreshExpiration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to create a token")))
		return
	}

	// create the session to the database
	err = au.DBStore.SessionRepository.Create(ctx, &models.Session{
		ID:                  payloadrefresh.Id.String(),
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

func (au *appU) RefreshTokenController(ctx *gin.Context) {
	var req RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResp(err))
		return
	}

	// verify the refresh token
	payload, err := au.TokenCreator.Verify(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResp(errors.New("invalid token")))
		return
	}

	// get the session from the database
	session, err := au.DBStore.SessionRepository.FindByID(ctx, payload.Id.String())

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResp(err))
		return
	}

	// some checks
	if session.UserID != payload.UserId {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResp(errors.New("Session does not belong to the user")))
		return
	}

	// hmm, is that right ?
	if time.Now().After(session.RefreshTokenExpires) {
		err = fmt.Errorf("Session has been expired before!")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(err))
	}

	// create a new access token
	token, payload, err := au.TokenCreator.Create(session.UserID, au.AppConfig.TokenAccessExpiration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to create a token")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Token has been refreshed successfully",
		"access_token":  token,
		"refresh_token": session.RefreshToken,
	})
}
