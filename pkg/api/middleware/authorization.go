package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ayehia0/org/pkg/token"
	"github.com/ayehia0/org/pkg/utils"
	"github.com/gin-gonic/gin"
)

var (
	authorizationHeaderKey = "authorization"
	authorizationType      = "bearer"

	// to be able to store the token in the context, so we can access it later
	AuthPayloadKey = "authorization_payload_ctx"
)

func AuthMiddleware(tokenCreator token.TokenCreator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// implement the authentication here
		// check the header : authentication
		authHeader := ctx.Request.Header.Get(authorizationHeaderKey)
		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				utils.ErrorResp(errors.New("Authentication header is empty")),
			)
			return
		}
		// get the token
		authFields := strings.Fields(authHeader)

		if len(authFields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				utils.ErrorResp(errors.New("Invalid authentication header format")),
			)
			return
		}

		// verifiy the token
		authType := strings.ToLower(authFields[0])
		if authType != authorizationType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				utils.ErrorResp(fmt.Errorf("Unspported authorization type %s", authType)),
			)
			return
		}

		payload, err := tokenCreator.Verify(authFields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResp(err))
			return
		}
		ctx.Set(AuthPayloadKey, payload)
		ctx.Next()
	}
}
