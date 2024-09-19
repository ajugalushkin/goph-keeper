package grpcapp

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	authv1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
	mocksAuth "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/auth/v1/mocks"
	mocksKeeper "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/keeper/v1/mocks"
	"github.com/ajugalushkin/goph-keeper/server/internal/services"
)

// Initializes App struct correctly with provided parameters
func TestNewAppInitialization(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	authService := mocksAuth.NewAuth(t)
	keeperService := mocksKeeper.NewKeeper(t)

	jwtManager := services.NewJWTManager(log, "test_secret", time.Hour)

	address := "localhost:50051"

	app := New(log, authService, keeperService, jwtManager, address)

	if app.log != log {
		t.Errorf("expected log to be %v, got %v", log, app.log)
	}
	if app.gRPCServer == nil {
		t.Error("expected gRPCServer to be initialized, got nil")
	}
	if app.address != address {
		t.Errorf("expected address to be %s, got %s", address, app.address)
	}
}

// Returns a list containing the full method names for RegisterV1 and LoginV1
func TestAccessibleMethodsReturnsCorrectMethodNames(t *testing.T) {
	expected := []string{
		authv1.AuthServiceV1_RegisterV1_FullMethodName,
		authv1.AuthServiceV1_LoginV1_FullMethodName,
	}

	result := accessibleMethods()

	assert.Equal(t, expected, result)
}

// MustRun successfully calls Run without errors
func TestMustRun_Success(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	authService := mocksAuth.NewAuth(t)
	keeperService := mocksKeeper.NewKeeper(t)

	jwtManager := services.NewJWTManager(log, "secret", 15*time.Minute)
	app := New(log, authService, keeperService, jwtManager, ":50051")

	go func() {
		time.Sleep(1 * time.Second)
		app.Stop()
	}()

	app.MustRun()
}

// MustRun handles the scenario where Run returns an error
func TestMustRun_RunError(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	authService := mocksAuth.NewAuth(t)
	keeperService := mocksKeeper.NewKeeper(t)

	jwtManager := services.NewJWTManager(log, "secret", 15*time.Minute)
	app := New(log, authService, keeperService, jwtManager, ":invalid")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic but did not get one")
		}
	}()

	app.MustRun()
}

// Returns an error if the provided address is invalid or already in use
func TestRun_InvalidAddress(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	authService := mocksAuth.NewAuth(t)
	keeperService := mocksKeeper.NewKeeper(t)

	jwtManager := services.NewJWTManager(log, "secret", 15*time.Minute)
	app := New(log, authService, keeperService, jwtManager, "invalid_address")

	err := app.Run()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
