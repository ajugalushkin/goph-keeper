package app

import (
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/server/config"
	grpcapp "github.com/ajugalushkin/goph-keeper/server/internal/app/grpc"
	"github.com/ajugalushkin/goph-keeper/server/internal/lib/jwt"
	"github.com/ajugalushkin/goph-keeper/server/internal/services/auth"
	"github.com/ajugalushkin/goph-keeper/server/internal/services/keeper"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage/postgres"
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

	jwtManager := jwt.NewJWTManager(log, cfg.Token.TokenSecret, cfg.Token.TokenTTL)

	serviceAuth := auth.New(log, storage, storage, jwtManager)
	serviceKeeper := keeper.New(log, storage, storage)

	grpcApp := grpcapp.New(
		log,
		serviceAuth,
		serviceKeeper,
		jwtManager,
		cfg.GRPC.ServerAddress,
	)

	return &App{
		GRPCSrv: grpcApp,
	}
}
