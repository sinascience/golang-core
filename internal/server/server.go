package server

import (
	"log/slog"
	"os"
	"sync"
	"time"
	"venturo-core/configs"
	"venturo-core/internal/database"

	"github.com/gofiber/fiber/v2"
)

// NewServer creates and configures a new Fiber application.
func NewServer() (*fiber.App, *sync.WaitGroup) {
	config, err := configs.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	database.ConnectDB(&config)

	app := fiber.New()

	var wg sync.WaitGroup

	registerRoutes(app, database.DB, &config, &wg)

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
