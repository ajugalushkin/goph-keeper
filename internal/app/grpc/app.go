package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
)

type App struct {
	GRPCSrv *grpc.Server
}

func New(
	log *slog.Logger,
	grpcPort int,
	StoragePath string,
	tokenTTL time.Duration,
) *App {
}
func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op)).
		Info("grpc server is stopping", slog.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
