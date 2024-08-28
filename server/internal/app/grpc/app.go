package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	keeperv1 "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/keeper/v1"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	address    string
}

// New функция создания экземпляра приложения
func New(
	log *slog.Logger,
	authService keeperv1.Auth,
	Address string,
) *App {
	gRPCServer := grpc.NewServer()
	keeperv1.Register(gRPCServer, authService)

	// Register reflection service on gRPC server.
	reflection.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		address:    Address,
	}
}

// MustRun метод для запуска приложения, при возникновении ошибки паникуем
func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.String("address", a.address))

	listener, err := net.Listen("tcp", a.address)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op)).
		Info("grpc server is stopping", slog.String("address", a.address))
	a.gRPCServer.GracefulStop()
}
