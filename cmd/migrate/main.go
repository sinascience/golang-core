package main

import (
	"log/slog"
	"os"
	"venturo-core/configs"
	"venturo-core/internal/database"
)

func main() {
	slog.Info("Migration tool started")

	config, err := configs.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	database.ConnectDB(&config)

	if len(os.Args) < 2 {
		slog.Error("Please provide an argument: up, down, or fresh")
		os.Exit(1)
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
		slog.Error("Unknown command", "command", command)
		os.Exit(1)
	}
}
