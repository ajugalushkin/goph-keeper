package app

import (
	"log/slog"

	"google.golang.org/grpc"

	authgrpcv1 "github.com/ajugalushkin/goph-keeper/internal/handlers/grpc/auth/v1"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	port int,
) *App {
	grpcServer := grpc.NewServer()
	authgrpcv1.Register(grpcServer)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}
