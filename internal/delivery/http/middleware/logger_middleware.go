package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/infrastructure/logger"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		requestID := c.GetString(RequestIDKey)

		fields := []any{
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
			"request_id", requestID,
		}

		if query != "" {
			fields = append(fields, "query", query)
		}

		if len(c.Errors) > 0 {
			fields = append(fields, "errors", c.Errors.String())
		}

		if status >= 500 {
			logger.Error("HTTP request completed with server error", fields...)
		} else if status >= 400 {
			logger.Warn("HTTP request completed with client error", fields...)
		} else {
			logger.Info("HTTP request completed", fields...)
		}
	}
}
