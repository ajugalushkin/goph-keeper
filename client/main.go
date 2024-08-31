package main

import (
	"github.com/ajugalushkin/goph-keeper/client/internal/cli"
	"log/slog"
	"os"

	"github.com/ajugalushkin/goph-keeper/client/config"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info(
		"starting client",
		slog.String("env", cfg.Env),
		slog.String("build version", buildVersion),
		slog.String("build date", buildDate),
	)

	//_, err := client.New(context.Background(), log, cfg.Client.Address, cfg.Client.Timeout, cfg.Client.RetriesCount)
	//if err != nil {
	//	log.Info("failed to initialize client")
	//}

	//tui.MustRun()

	cli.Execute()

	//// Graceful shutdown
	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//
	//sign := <-stop
	//log.Info("stopping application", slog.String("signal", sign.String()))
	//
	////application.GRPCSrv.Stop()
	//log.Info("application stopped")
}

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
