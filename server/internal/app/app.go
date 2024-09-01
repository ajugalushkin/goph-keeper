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

	userStorage, err := postgres.NewUserStorage(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	vaultStorage, err := postgres.NewVaultStorage(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	jwtManager := jwt.NewJWTManager(log, cfg.Token.TokenSecret, cfg.Token.TokenTTL)

	serviceAuth := auth.New(log, userStorage, userStorage, jwtManager)
	serviceKeeper := keeper.New(log, vaultStorage, vaultStorage)

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
