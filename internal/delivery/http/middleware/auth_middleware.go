package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/delivery/http/response"
	"github.com/tms/tyre/internal/infrastructure/jwt"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	ContextUserID       = "user_id"
	ContextEmail        = "email"
	ContextRole         = "role"
	ContextTenantID     = "tenant_id"
	ContextClaims       = "claims"
)

func GetClaims(c *gin.Context) (*jwt.JWTClaims, bool) {
	claimsVal, exists := c.Get(ContextClaims)
	if !exists {
		return nil, false
	}
	claims, ok := claimsVal.(*jwt.JWTClaims)
	return claims, ok
}

func Auth(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Unauthorized(c, "Invalid authorization format. Use: Bearer <token>")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				response.Error(c, http.StatusUnauthorized, "Token has expired", nil)
			} else {
				response.Error(c, http.StatusUnauthorized, "Invalid token", nil)
			}
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextRole, claims.Role)
		c.Set(ContextTenantID, claims.TenantID)
		c.Set(ContextClaims, claims)

		c.Next()
	}
}
