package main

import (
	"github.com/ajugalushkin/goph-keeper/server/config"
	"github.com/ajugalushkin/goph-keeper/server/internal/app"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

// main is the entry point of the Goph-Keeper server application.
// It initializes the application configuration, sets up a logger based on the environment,
// creates a new application instance, runs the gRPC server, waits for OS signals for graceful shutdown,
// and logs the application's start and stop events.
func main() {
	cfg := config.MustLoad()
	// Set up a logger based on the application environment.
	log := setupLogger(cfg.Env)

	application := initApp(log, cfg)

	run(application)

	waitSignal(log, application)
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

// run starts the gRPC server in a separate goroutine.
//
// The function takes an instance of the application, which contains the gRPC server.
// It runs the gRPC server in a separate goroutine using the MustRun method of the gRPC server.
// This allows the main function to continue executing without waiting for the gRPC server to finish.
//
// Parameters:
// - application: A pointer to an instance of the application struct, which contains the gRPC server.
func run(application *app.App) {
	// Run the gRPC server in a separate goroutine.
	go application.GRPCSrv.MustRun()
}

// waitSignal listens for OS signals to gracefully shut down the application.
// It sets up a channel to receive SIGINT, SIGTERM, and SIGQUIT signals, waits for a signal,
// logs the received signal, stops the gRPC server, and logs the successful shutdown of the application.
//
// Parameters:
// - log: A pointer to a slog.Logger instance for logging events.
// - application: A pointer to an instance of the application struct, which contains the gRPC server.
func waitSignal(log *slog.Logger, application *app.App) {
	// Set up a channel to receive OS signals for graceful shutdown.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Wait for a signal to waitSignal the application.
	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	// Stop the gRPC server.
	application.GRPCSrv.Stop()

	// Log the successful waitSignal of the application.
	log.Info("application stopped")
}

// initApp initializes and returns a new application instance with the provided logger and configuration.
// It logs the start of the application with its environment, build version, and build date.
//
// Parameters:
// - log: A pointer to a slog.Logger instance for logging events.
// - cfg: A pointer to a config2.Config instance containing the application's configuration.
//
// Returns:
// - A pointer to a new instance of the app.App struct, which contains the gRPC server.
func initApp(log *slog.Logger, cfg *config.Config) *app.App {
	// Log the start of the application with its environment, build version, and build date.
	log.Info(
		"starting application",
		slog.String("env", cfg.Env),
		slog.String("build version", buildVersion),
		slog.String("build date", buildDate),
	)

	// Create a new application instance with the configured logger and configuration.
	return app.New(log, cfg)
}
