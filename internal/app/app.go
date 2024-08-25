package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/ajugalushkin/goph-keeper/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcAddress string,
	tokenTTL time.Duration,
) *App {
	grpcApp := grpcapp.New(log, grpcAddress)
	return &App{
		GRPCSrv: grpcApp,
	}
}
