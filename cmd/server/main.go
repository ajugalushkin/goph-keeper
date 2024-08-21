package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ajugalushkin/goph-keeper/internal/app"
	"github.com/ajugalushkin/goph-keeper/internal/config"
)

const (
	envDev    = "debug"
	envProd   = "info"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application")

	application := app.New(log,00,"",cfg.TokenTTL)

	application.GRPCSrv.MustRun()
}

func setupLogger(Env string) *slog.Logger {
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

	return log
}
