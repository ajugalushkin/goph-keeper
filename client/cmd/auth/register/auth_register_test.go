package register

import (
	"log/slog"
	"testing"

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
func TestNewCommand_ValidLoggerAndNilClient(t *testing.T) {
	t.Parallel()

	// Arrange
	var logger *slog.Logger
	var client app.AuthClient

	// Act
	cmd := NewCommand(logger, client)

	// Assert
	if cmd == nil {
		t.Error("Expected NewCommand to return a non-nil command")
	}
	if register == nil {
		t.Error("Expected register to be initialized")
	}
	if register.client != nil {
		t.Error("Expected register.client to be nil")
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
