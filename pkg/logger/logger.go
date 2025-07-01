package logger

import (
	"log/slog"
	"os"
)

// InitLogger sets up our structured, JSON-based logger.
func InitLogger() {
	// Create a new JSON handler that writes to standard output.
	// We can set a log level, e.g., slog.LevelDebug to see all logs.
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// Set this handler as the default logger for the entire application.
	slog.SetDefault(slog.New(handler))
}
