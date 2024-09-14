package register

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
)

func TestNewCommand_NilLoggerNilClient(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("TestNewCommand_NilLoggerNilClient panicked: %v", r)
		}
	}()

	var nilLogger *slog.Logger
	var nilClient app.AuthClient

	cmd := NewCommand(nilLogger, nilClient)

	if cmd == nil {
		t.Error("Expected NewCommand to return a non-nil command")
	}
}

func TestNewCommand_NilLoggerValidClient(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("TestNewCommand_NilLoggerValidClient panicked: %v", r)
		}
	}()

	var nilLogger *slog.Logger
	var validClient app.AuthClient = mocks.NewAuthClient(t)

	cmd := NewCommand(nilLogger, validClient)

	if cmd == nil {
		t.Error("Expected NewCommand to return a non-nil command")
	}
	if register == nil {
		t.Error("Expected register to be initialized")
	}
	if register.client == nil {
		t.Error("Expected register.client to be initialized")
	}
}

func TestNewCommand_ValidLoggerAndValidClient(t *testing.T) {
	t.Parallel()

	// Arrange
	var logger *slog.Logger
	var client app.AuthClient = mocks.NewAuthClient(t)

	// Act
	cmd := NewCommand(logger, client)

	// Assert
	if cmd == nil {
		t.Error("Expected NewCommand to return a non-nil command")
	}
	if register == nil {
		t.Error("Expected register to be initialized")
	}
	if register.client == nil {
		t.Error("Expected register.client to be initialized")
	}
}

// Successfully retrieves email and password from command-line flags
func TestRegisterCmdRun_Success(t *testing.T) {
	// Arrange
	mockAuthClient := mocks.NewAuthClient(t)
	mockAuthClient.On("Register", mock.Anything, "test@example.com", "password123").Return(nil)

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cmd := NewCommand(log, mockAuthClient)
	cmd.Flags().Set("email", "test@example.com")
	cmd.Flags().Set("password", "password123")

	// Act
	err := registerCmdRun(cmd, []string{})

	// Assert
	assert.NoError(t, err)
}

func TestRegisterCmdRun_EmptyEmail(t *testing.T) {
	// Arrange
	mockAuthClient := mocks.NewAuthClient(t)
	mockAuthClient.On(
		"Register",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil).Maybe()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cmd := NewCommand(log, mockAuthClient)
	cmd.Flags().Set("email", "")
	cmd.Flags().Set("password", "password123")

	// Act
	err := registerCmdRun(cmd, []string{})

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "email is required")
}

func TestRegisterCmdRun_EmptyPassword(t *testing.T) {
	// Arrange
	mockAuthClient := mocks.NewAuthClient(t)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cmd := NewCommand(log, mockAuthClient)
	cmd.Flags().Set("email", "test@example.com")
	cmd.Flags().Set("password", "")

	// Act
	err := registerCmdRun(cmd, []string{})

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "password is required")
}
