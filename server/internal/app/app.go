package app

import (
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/server/config"
	"github.com/ajugalushkin/goph-keeper/server/internal/app/grpc"
	"github.com/ajugalushkin/goph-keeper/server/internal/services"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage/minio"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage/postgres"
)

type App struct {
	GRPCSrv grpcapp.GrpcServer
}

// New initializes and returns a new instance of the App struct.
// The function takes two parameters:
// - log: A pointer to a slog.Logger instance for logging.
// - cfg: A pointer to a config.Config instance containing application configuration.
//
// The function initializes various storage and service components based on the provided configuration.
// It creates instances of UserStorage, VaultStorage, MinioStorage, JWTManager, AuthService, and KeeperService.
// Then, it creates a new instance of GrpcServer using the initialized services and configuration.
//
// Finally, it returns a new App instance containing the initialized GrpcServer.
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	jwtManager := services.NewJWTManager(log, cfg.Token.Secret, cfg.Token.TTL)

	serviceAuth := initAuthService(log, cfg, jwtManager)

	serviceKeeper := initKeeperService(log, cfg)

	grpcApp := grpcapp.New(
		log,
		serviceAuth,
		serviceKeeper,
		jwtManager,
		cfg.GRPC.Address,
	)

	return &App{
		GRPCSrv: grpcApp,
	}
}

func initAuthService(
	log *slog.Logger,
	cfg *config.Config,
	jwtManager *services.JWTManager,
) *services.Auth {
	userStorage, err := postgres.NewUserStorage(cfg.Storage.Path)
	if err != nil {
		panic(err)
	}

	if jwtManager == nil {
		panic("JWT manager is nil")
	}

	return services.NewAuthService(
		log,
		userStorage,
		userStorage,
		jwtManager,
	)
}

func initKeeperService(log *slog.Logger, cfg *config.Config) *services.Keeper {
	vaultStorage, err := postgres.NewVaultStorage(cfg.Storage.Path)
	if err != nil {
		panic(err)
	}

	minioStorage, err := minio.NewMinioStorage(cfg.Minio)
	if err != nil {
		panic(err)
	}
	return services.NewKeeperService(
		log,
		vaultStorage,
		vaultStorage,
		minioStorage,
		minioStorage,
	)
}
