package logger

import (
	"log/slog"
	"os"
	"sync"

	"github.com/ajugalushkin/goph-keeper/client/config"
)

type LogInstance struct {
	Log *slog.Logger
}

var (
	log  *LogInstance
	once sync.Once
)

func GetInstance() *LogInstance {
	once.Do(
		func() {
			cfg := config.GetInstance().Config
			log = &LogInstance{
				Log: setupLogger(cfg.Env),
			}
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
