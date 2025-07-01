package server

import (
	"log"
	"sync"
	"time"
	"venturo-core/configs"
	"venturo-core/internal/database"

	"github.com/gofiber/fiber/v2"
)

// NewServer creates and configures a new Fiber application.
func NewServer() (*fiber.App, *sync.WaitGroup) {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("🔥 Failed to load configuration: %v", err)
	}

	// Connect to the database
	database.ConnectDB(&config)

	// Initialize Fiber app
	app := fiber.New()

	// Create the single WaitGroup instance here.
	var wg sync.WaitGroup

	// Pass the WaitGroup down to the router setup.
	registerRoutes(app, database.DB, &config, &wg)

	return app, &wg
}

func GracefulShutdown(app *fiber.App, wg *sync.WaitGroup) {
	log.Println("Gracefully shutting down...")
	log.Println("Waiting for background processes to finish...")
	wg.Wait() // Wait on the single, shared WaitGroup
	log.Println("All background processes finished.")

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
