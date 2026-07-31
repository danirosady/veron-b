package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			allowed := false

			// Check exact match first
			for _, o := range allowedOrigins {
				if origin == o {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					allowed = true
					break
				}
			}

			// In development, allow localhost origins
			if !allowed {
				if strings.HasPrefix(origin, "http://localhost:") ||
					strings.HasPrefix(origin, "http://127.0.0.1:") {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
