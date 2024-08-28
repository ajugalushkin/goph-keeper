package app

import (
	"log/slog"

	grpcapp "github.com/ajugalushkin/goph-keeper/internal/app/grpc"
	"github.com/ajugalushkin/goph-keeper/internal/config"
	"github.com/ajugalushkin/goph-keeper/internal/services/auth"
	"github.com/ajugalushkin/goph-keeper/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {

	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	serviceAuth := auth.New(log, storage, storage, cfg.TokenTTL)

	grpcApp := grpcapp.New(log, serviceAuth, cfg.GRPC.ServerAddress)
	return &App{
		GRPCSrv: grpcApp,
	}
}
