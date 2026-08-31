package middlewares

import (
	"gotask-api/constants"
	"gotask-api/datatransfers"
	"gotask-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware, for authentication user JWT and attach user info to context

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// check the authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			datatransfers.ErrorRes(c, http.StatusUnauthorized, constants.UNAUTHORIZED, "missing authorization header")
			c.Abort()
			return
		}

		// step 2, check the beare
		bearerCheck := strings.SplitN(authHeader, " ", 2)
		if len(bearerCheck) != 2 || bearerCheck[0] != "Bearer" {
			datatransfers.ErrorRes(c, http.StatusUnauthorized, constants.UNAUTHORIZED, "invalid authorization format, expected 'Bearer <token>'")
			c.Abort()
			return
		}
		// get the real token
		tokenString := bearerCheck[1]

		// Step 3 — validate the JWT
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			datatransfers.ErrorRes(c, http.StatusUnauthorized, constants.UNAUTHORIZED, "invalid or expired token: "+err.Error())
			c.Abort()
			return
		}

		// Step 4 — attach user info to context for handlers to use
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		// Step 5 — continue to the next handler
		c.Next()

	}

}
