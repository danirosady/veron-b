package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/infrastructure/logger"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString(RequestIDKey)

				logger.Error("Panic recovered",
					"error", fmt.Sprintf("%v", err),
					"request_id", requestID,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"stack", string(debug.Stack()),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Internal server error",
					"error":   "An unexpected error occurred",
				})
			}
		}()

		c.Next()
	}
}
