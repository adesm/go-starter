package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strings"

	"boilerplate/internal/shared/response"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := make([]byte, 2048)
				stack = stack[:runtime.Stack(stack, false)]

				// Extract file and line from stack
				stackLines := strings.Split(string(stack), "\n")
				location := "unknown"
				if len(stackLines) > 3 {
					// Index 3 usually contains the caller info in a panic scenario
					location = strings.TrimSpace(stackLines[3])
				}

				// Log the error deeply for developers
				slog.Error("PANIC RECOVERED",
					slog.Any("error", err),
					slog.String("path", c.Request.URL.Path),
					slog.String("method", c.Request.Method),
					slog.String("location", location),
					slog.String("stack", string(stack)),
				)

				// Also print to console for development visual debug
				fmt.Printf("\n[RECOVERED PANIC]\nError: %v\nLocation: %s\n", err, location)

				// Send clean response to client
				response.JSONError(c, http.StatusInternalServerError, fmt.Errorf("%v", err))
				c.Abort()
			}
		}()
		c.Next()
	}
}
