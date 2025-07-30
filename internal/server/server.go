package server

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"
	"venturo-core/configs"
	"venturo-core/internal/database"
	"venturo-core/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// NewServer creates and configures a new Fiber application.
func NewServer() (*fiber.App, *sync.WaitGroup) {
	config, err := configs.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	db := database.ConnectDB(&config)

	app := fiber.New()
	
	// Add middleware to set client IP in context
	app.Use(func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.Context(), "client_ip", c.IP())
		ctx = context.WithValue(ctx, "request_time", time.Now())
		c.SetUserContext(ctx)
		return c.Next()
	})
	
	// Add logging middleware
	app.Use(middleware.LoggingMiddleware())
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.CORSAllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	var wg sync.WaitGroup

	registerRoutes(app, db, &config, &wg)

	return app, &wg
}

func GracefulShutdown(app *fiber.App, wg *sync.WaitGroup) {
	slog.Info("Gracefully shutting down...")
	slog.Info("Waiting for background processes to finish...")
	wg.Wait()
	slog.Info("All background processes finished.")

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Server gracefully stopped.")
}
