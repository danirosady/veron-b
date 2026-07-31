package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/delivery/http/response"
)

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextRole)
		if !exists {
			response.Forbidden(c, "Role not found in context")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			response.Forbidden(c, "Invalid role type in context")
			c.Abort()
			return
		}

		if !slices.Contains(allowedRoles, roleStr) {
			response.Error(c, http.StatusForbidden, "Insufficient permissions", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
