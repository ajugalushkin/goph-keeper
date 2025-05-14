package auth

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	authv1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
	"github.com/ajugalushkin/goph-keeper/mocks"
)

func TestRegister_EmptyPassword(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAuthService := mocks.NewAuthServiceV1Client(t)
	mockAuthService.On("RegisterV1",
		ctx,
		&authv1.RegisterRequestV1{Email: "user@example.com", Password: ""},
	).Return(
		nil,
		status.Error(codes.InvalidArgument, "invalid password"),
	)

	client := &Client{api: mockAuthService}

	// Act
	err := client.Register(ctx, "user@example.com", "")

	// Assert
	require.Error(t, err)
	require.EqualError(t, err, "client.auth.Register: rpc error: code = InvalidArgument desc = invalid password")
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAuthService := mocks.NewAuthServiceV1Client(t)
	mockAuthService.On("RegisterV1",
		ctx,
		&authv1.RegisterRequestV1{Email: "existing@email.com", Password: "password"},
	).Return(
		nil,
		status.Error(codes.AlreadyExists, "email already registered"),
	)

	client := &Client{api: mockAuthService}

	// Act
	err := client.Register(ctx, "existing@email.com", "password")

	// Assert
	require.Error(t, err)
	require.EqualError(t, err, "client.auth.Register: rpc error: code = AlreadyExists desc = email already registered")
}

func TestRegister_EmptyEmail(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAuthService := mocks.NewAuthServiceV1Client(t)
	mockAuthService.On("RegisterV1",
		ctx,
		&authv1.RegisterRequestV1{Email: "", Password: "password"},
	).Return(
		nil,
		status.Error(codes.InvalidArgument, "invalid email address"),
	)

	client := &Client{api: mockAuthService}

	// Act
	err := client.Register(ctx, "", "password")

	// Assert
	require.Error(t, err)
	require.EqualError(t, err, "client.auth.Register: rpc error: code = InvalidArgument desc = invalid email address")
}

func TestRegister_ValidEmail(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAuthService := mocks.NewAuthServiceV1Client(t)
	mockAuthService.On("RegisterV1",
		ctx,
		&authv1.RegisterRequestV1{Email: "valid@email.com", Password: "password"},
	).Return(
		&authv1.RegisterResponseV1{UserId: 42},
		nil,
	)

	client := &Client{api: mockAuthService}

	// Act
	err := client.Register(ctx, "valid@email.com", "password")

	// Assert
	require.NoError(t, err)
	mockAuthService.AssertExpectations(t)
}

func TestLogin_EmptyEmail(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAuthService := mocks.NewAuthServiceV1Client(t)
	mockAuthService.On("LoginV1",
		ctx,
		&authv1.LoginRequestV1{Email: "", Password: "password"},
	).Return(
		nil,
		status.Error(codes.InvalidArgument, "invalid email address"),
	)

	client := &Client{api: mockAuthService}

	// Act
	token, err := client.Login(ctx, "", "password")

	// Assert
	require.Error(t, err)
	require.Empty(t, token)
	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = invalid email address")
}
func TestLogin_EmptyPassword(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAuthService := mocks.NewAuthServiceV1Client(t)
	mockAuthService.On("LoginV1",
		ctx,
		&authv1.LoginRequestV1{Email: "user@example.com", Password: ""},
	).Return(
		nil,
		status.Error(codes.InvalidArgument, "invalid password"),
	)

	client := &Client{api: mockAuthService}

	// Act
	token, err := client.Login(ctx, "user@example.com", "")

	// Assert
	require.Error(t, err)
	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = invalid password")
	require.Empty(t, token)
}
func TestLogin_EmailNotRegistered(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAuthService := mocks.NewAuthServiceV1Client(t)
	mockAuthService.On("LoginV1",
		ctx,
		&authv1.LoginRequestV1{Email: "not_registered@example.com", Password: "password"},
	).Return(
		nil,
		status.Error(codes.NotFound, "email not registered"),
	)

	client := &Client{api: mockAuthService}

	// Act
	token, err := client.Login(ctx, "not_registered@example.com", "password")

	// Assert
	require.Error(t, err)
	require.EqualError(t, err, "rpc error: code = NotFound desc = email not registered")
	require.Empty(t, token)
}

func TestGetAuthConnection_ConnectionError(t *testing.T) {
	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})), &config.Config{Env: "dev"})
	log := logger.GetLogger()

	cfg := config.Client{
		Address: "invalid_address",
		Retries: 3,
		Timeout: 5 * time.Second,
	}

	// Act
	GetAuthConnection(log, cfg)
}
