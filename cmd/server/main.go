package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"venturo-core/internal/server"
	"venturo-core/pkg/logger"
)

func main() {
	// Get the app and the shared WaitGroup from our server setup
	app, wg := server.NewServer()
	logger.InitLogger()

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Server is starting", "port", 3000)
		if err := app.Listen(":3000"); err != nil {
			slog.Error("Server failed to start", "error", err)
		}
	}()

	// Block until a signal is received
	<-quit

	// Trigger the graceful shutdown, passing the shared WaitGroup
	server.GracefulShutdown(app, wg)
}
