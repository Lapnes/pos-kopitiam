package middleware

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.ErrorResponse("Role not found in context", "FORBIDDEN", nil))
			return
		}

		userRole := role.(string)
		isAllowed := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.ErrorResponse("Insufficient permissions", "FORBIDDEN", nil))
			return
		}

		c.Next()
	}
}
