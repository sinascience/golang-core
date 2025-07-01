package main

import (
	"log"
	"os"
	"venturo-core/configs"
	"venturo-core/internal/database"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("🔥 Failed to load configuration: %v", err)
	}

	// Connect to the database
	database.ConnectDB(&config)

	// Check for command-line arguments
	if len(os.Args) < 2 {
		log.Fatal("Please provide an argument: up, down, or fresh")
		return
	}

	command := os.Args[1]

	switch command {
	case "up":
		database.MigrateUp()
	case "down":
		database.MigrateDown()
	case "fresh":
		database.Drop()
		database.MigrateUp()
	default:
		log.Fatalf("Unknown command: %s. Please use 'up', 'down', or 'fresh'.", command)
	}
}
