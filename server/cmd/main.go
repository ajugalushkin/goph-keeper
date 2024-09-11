package main

import (
	"github.com/ajugalushkin/goph-keeper/server/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ajugalushkin/goph-keeper/server/internal/app"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

// main is the entry point of the Goph-Keeper server application.
// It initializes the application configuration, sets up a logger, starts the gRPC server,
// and handles graceful shutdowns.
func main() {
    cfg := config.MustLoad()

    // Set up a logger based on the application environment.
    log := setupLogger(cfg.Env)

    // Log the start of the application with its environment, build version, and build date.
    log.Info(
        "starting application",
        slog.String("env", cfg.Env),
        slog.String("build version", buildVersion),
        slog.String("build date", buildDate),
    )

    // Create a new application instance with the configured logger and configuration.
    application := app.New(log, cfg)

    // Run the gRPC server in a separate goroutine.
    go application.GRPCSrv.MustRun()

    // Set up a channel to receive OS signals for graceful shutdown.
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

    // Wait for a signal to stop the application.
    sign := <-stop
    log.Info("stopping application", slog.String("signal", sign.String()))

    // Stop the gRPC server.
    application.GRPCSrv.Stop()

    // Log the successful stop of the application.
    log.Info("application stopped")
}

// setupLogger initializes and configures a slog.Logger based on the provided environment.
// It creates a logger with different log levels based on the environment:
// - In development mode (envDev), it sets the log level to Debug.
// - In production mode (envProd), it sets the log level to Info.
//
// Parameters:
// - Env (string): The environment in which the application is running. It can be either "dev" or "prod".
//
// Returns:
// - *slog.Logger: A configured slog.Logger instance.
func setupLogger(Env string) *slog.Logger {
    const (
        envDev  = "dev"
        envProd = "prod"
    )

    var log *slog.Logger

    switch Env {
    case envDev:
        log = slog.New(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    case envProd:
        log = slog.New(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
        )
    }
    return log
}
