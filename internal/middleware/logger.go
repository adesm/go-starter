package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		slog.Info("http request",
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.Int("status", status),
			slog.Duration("duration", duration),
		)
	}
}
