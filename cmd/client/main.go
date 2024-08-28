package main

import (
	"context"
	"log/slog"
	"os"

	authclient "github.com/ajugalushkin/goph-keeper/internal/app/client/auth"
	"github.com/ajugalushkin/goph-keeper/internal/config"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	_, err := authclient.New(context.Background(), log, cfg.Client.Address, cfg.Client.Timeout, cfg.Client.RetriesCount)
	if err != nil {
		return
	}
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
	}
	return log
}
