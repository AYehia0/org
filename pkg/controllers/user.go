package controllers

import (
	"errors"
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

	// revoke refresh token
	RevokeRefreshTokenController(ctx *gin.Context)
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

	// Save the session to the database
	session := &models.Session{
		ID:                  payloadrefresh.Id.String(),
		UserID:              user.ID,
		AccessToken:         userAccessToken,
		RefreshToken:        userRefreshToken,
		AccessTokenExpires:  payloadAccess.ExpiredAt,
		RefreshTokenExpires: payloadrefresh.ExpiredAt,
	}
	err = au.DBStore.SessionRepository.Create(ctx, session)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to create a session")))
		return
	}

	// Also save the refresh token to the redis database
	err = au.RDBStore.SessionRepository.CreateSession(ctx, session)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to save session to the redis")))
		return
	}

	ctx.JSON(http.StatusOK, returnRefreshTokenResponse(session.RefreshToken, userAccessToken))
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

	// check if the token is valid in the redis database first
	sessionInRedis, err := au.RDBStore.SessionRepository.GetSessionByID(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResp(err))
		return
	}

	// if the token is found
	if sessionInRedis != nil {
		err = isSessionValid(sessionInRedis)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResp(err))
			return
		}

		// create a new access token
		token, _, err := au.TokenCreator.Create(sessionInRedis.ID, au.AppConfig.TokenAccessExpiration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to create a token")))
			return
		}
		ctx.JSON(http.StatusOK, returnRefreshTokenResponse(sessionInRedis.RefreshToken, token))
		return
	}

	// get the session from the database
	session, err := au.DBStore.SessionRepository.FindByID(ctx, payload.Id.String())
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResp(err))
		return
	}

	err = isSessionValid(session)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResp(err))
		return
	}
	// create a new access token
	token, payload, err := au.TokenCreator.Create(session.ID, au.AppConfig.TokenAccessExpiration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(errors.New("failed to create a token")))
		return
	}
	ctx.JSON(http.StatusOK, returnRefreshTokenResponse(session.RefreshToken, token))
}

func (au *appU) RevokeRefreshTokenController(ctx *gin.Context) {
	var req RevokeRefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResp(err))
		return
	}

	// delete the refresh token from the redis database
	err := au.RDBStore.SessionRepository.DeleteSession(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResp(err))
		return
	}

	// this is optional to delete the token from the database or now
	// TODO: delete the session record from the database

	ctx.JSON(http.StatusOK, gin.H{"message": "Token has been revoked successfully"})
}

// helper function to check if the session is valid
func isSessionValid(session *models.Session) error {
	// if the user isn't the owner of the session
	if session.UserID != session.UserID {
		return errors.New("Session does not belong to the user")
	}

	if time.Now().After(session.RefreshTokenExpires) {
		return errors.New("Session has been expired before!")
	}

	return nil
}

func returnRefreshTokenResponse(refreshToken, accessToken string) RefreshTokenResponse {
	return RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Token has been refreshed successfully",
	}
}
