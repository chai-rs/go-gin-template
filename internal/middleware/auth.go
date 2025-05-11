package middleware

import (
	"net/http"

	"github.com/0xanonydxck/simple-bookstore/infrastructure/auth"
	"github.com/0xanonydxck/simple-bookstore/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if the user is authenticated.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			utils.ResponseErrorWithStatus(c, http.StatusUnauthorized, "user hasn't logged in yet")
			c.Abort()
			return
		}

		c.Next()
	}
}

// Authorize checks if the user is authorized.
func Authorize(obj auth.AuthObject, act auth.AuthAction, enforcer auth.AuthEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			utils.ResponseErrorWithStatus(c, http.StatusUnauthorized, "user hasn't logged in yet")
			c.Abort()
			return
		}

		metadata, err := auth.ExtractTokenMetadata(c.Request)
		if err != nil {
			utils.ResponseErrorWithStatus(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		ok, err := enforcer.Enforce(metadata.UserID, obj, act)
		if err != nil {
			utils.ResponseErrorWithStatus(c, http.StatusUnauthorized, "error occurred while authorizing user")
			c.Abort()
			return
		}

		if !ok {
			utils.ResponseErrorWithStatus(c, http.StatusForbidden, "forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
