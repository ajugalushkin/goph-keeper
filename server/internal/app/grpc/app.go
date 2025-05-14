package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
	"github.com/ajugalushkin/goph-keeper/server/interceptors"
	authhandlerv1 "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/auth/v1"
	keeperhandlerv1 "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/services"
)

//go:generate mockery --name GrpcServer
type GrpcServer interface {
	MustRun()
	Run() error
	Stop()
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	address    string
}

// New creates a new gRPC application instance.
//
// The function initializes a new gRPC server with the provided services and configurations.
// It sets up an authentication interceptor to secure the gRPC endpoints and registers
// reflection service for debugging purposes.
//
// Parameters:
// - log: A pointer to a slog.Logger instance for logging.
// - authService: An instance of the authhandlerv1.Auth interface, representing the authentication service.
// - keeperService: An instance of the keeperhandlerv1.Keeper interface, representing the keeper service.
// - jwtManager: An instance of the services.TokenManager interface, used for managing JWT tokens.
// - Address: A string representing the address to listen on for incoming gRPC connections.
//
// Returns:
// - A pointer to a new App instance, containing the initialized gRPC server and other necessary components.
func New(
	log *slog.Logger,
	authService authhandlerv1.Auth,
	keeperService keeperhandlerv1.Keeper,
	jwtManager services.TokenManager,
	Address string,
) *App {

	interceptor := interceptors.NewAuthInterceptor(log, jwtManager, accessibleMethods())

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()))

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

// accessibleMethods returns a list of gRPC method names that are accessible without authentication.
// These methods are used to set up an authentication interceptor for securing the gRPC endpoints.
//
// Returns:
//   - A slice of strings, where each string represents a gRPC method name.
//     The method names are taken from the authv1.AuthServiceV1 service definition.
func accessibleMethods() []string {
	return []string{
		authv1.AuthServiceV1_RegisterV1_FullMethodName,
		authv1.AuthServiceV1_LoginV1_FullMethodName,
	}
}

// MustRun метод для запуска приложения, при возникновении ошибки паникуем
// MustRun starts the gRPC application and panics if an error occurs during startup.
// It is a convenience method that wraps the Run method and handles any errors by panicking.
//
// This method is intended to be used in main functions where a panic is acceptable.
// In production code, it is recommended to handle errors gracefully instead of panicking.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run starts the gRPC application and returns an error if any occurs during startup.
// It listens for incoming connections on the specified address and runs the gRPC server.
//
// Parameters:
// - a: A pointer to the App instance containing the gRPC server and other necessary components.
//
// Returns:
// - An error if any occurs during startup. If no error occurs, it returns nil.
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

// Stop gracefully stops the gRPC server and logs the shutdown process.
//
// This function is intended to be called when the application needs to gracefully shut down.
// It logs a message indicating that the server is stopping and then uses the GracefulStop method
// provided by the gRPC server to stop accepting new connections and finish processing existing ones.
//
// Parameters:
// - a: A pointer to the App instance containing the gRPC server and other necessary components.
//
// Returns:
// - This function does not return any value.
func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op)).
		Info("grpc server is stopping", slog.String("address", a.address))
	a.gRPCServer.GracefulStop()
}
