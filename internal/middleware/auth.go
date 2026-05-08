package middleware

import (
	"net/http"
	"strings"

	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("Authorization header missing", "UNAUTHORIZED", nil))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid token format", "UNAUTHORIZED", nil))
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateJWT(tokenString, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid or expired token", "UNAUTHORIZED", err.Error()))
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("branchID", claims.BranchID)

		c.Next()
	}
}
