package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// LoggingMiddleware creates a middleware for structured request logging
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request details
		duration := time.Since(start)

		// Determine log level based on status code
		status := c.Response().StatusCode()
		logLevel := slog.LevelInfo
		if status >= 400 && status < 500 {
			logLevel = slog.LevelWarn
		} else if status >= 500 {
			logLevel = slog.LevelError
		}

		// Extract user ID if available (from auth middleware)
		userID := "anonymous"
		if uid := c.Locals("current_user_id"); uid != nil {
			// Handle both string and UUID types
			switch v := uid.(type) {
			case string:
				userID = v
			case uuid.UUID:
				userID = v.String()
			default:
				userID = "unknown"
			}
		}

		slog.Log(c.Context(), logLevel, "HTTP Request",
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
			"user_id", userID,
			"response_size", len(c.Response().Body()),
		)

		return err
	}
}
