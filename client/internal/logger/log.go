package logger

import (
	"log/slog"
	"os"
	"sync"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
)

var (
	log  *slog.Logger
	once sync.Once
)

func InitLogger(newLog *slog.Logger) {
	log = newLog
}

func GetLogger() *slog.Logger {
	once.Do(
		func() {
			cfg := config.GetConfig()
			log = setupLogger(cfg.Env)
		})

	return log
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
