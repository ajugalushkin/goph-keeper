package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
	authhandlerv1 "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/auth/v1"
	keeperhandlerv1 "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/lib/jwt"
	"github.com/ajugalushkin/goph-keeper/server/internal/services/interceptors"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	address    string
}

// New функция создания экземпляра приложения
func New(
	log *slog.Logger,
	authService authhandlerv1.Auth,
	keeperService keeperhandlerv1.Keeper,
	jwtManager *jwt.JWTManager,
	Address string,
) *App {

	interceptor := interceptors.NewAuthInterceptor(log, jwtManager, accessibleMethods())

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()))

	keeperhandlerv1.Register(gRPCServer, keeperService)
	authhandlerv1.Register(gRPCServer, authService)

	// Register reflection service on gRPC server.
	reflection.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		address:    Address,
	}
}

func accessibleMethods() []string {
	return []string{
		authv1.AuthServiceV1_RegisterV1_FullMethodName,
		authv1.AuthServiceV1_LoginV1_FullMethodName,
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
