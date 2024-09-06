package app

import (
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/server/config"
	grpcapp "github.com/ajugalushkin/goph-keeper/server/internal/app/grpc"
	"github.com/ajugalushkin/goph-keeper/server/internal/services"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {

	userStorage, err := postgres.NewUserStorage(cfg.Storage.Path)
	if err != nil {
		panic(err)
	}

	vaultStorage, err := postgres.NewVaultStorage(cfg.Storage.Path)
	if err != nil {
		panic(err)
	}

	jwtManager := services.NewJWTManager(log, cfg.Token.Secret, cfg.Token.TTL)

	serviceAuth := services.NewAuthService(log, userStorage, userStorage, jwtManager)
	serviceKeeper := services.NewKeeperService(log, vaultStorage, vaultStorage)
	serviceMinio := services.NewMinioService(log, cfg.Minio)

	grpcApp := grpcapp.New(
		log,
		serviceAuth,
		serviceKeeper,
		serviceMinio,
		jwtManager,
		cfg.GRPC.Address,
	)

	return &App{
		GRPCSrv: grpcApp,
	}
}
