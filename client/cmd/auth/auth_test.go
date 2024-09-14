package auth

import (
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/ajugalushkin/goph-keeper/client/cmd/auth/login"
	"github.com/ajugalushkin/goph-keeper/client/cmd/auth/register"
	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
)

func TestNewCommand_ValidLoggerAndAuthClient(t *testing.T) {
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	expected := &cobra.Command{
		Use:   "auth",
		Short: "Manage user registration, authentication and authorization",
	}

	expected.AddCommand(login.NewCommand(mockLog, mockClient))
	expected.AddCommand(register.NewCommand(mockLog, mockClient))

	actual := NewCommand(mockLog, mockClient)

	assert.Equal(t, expected.Use, actual.Use)
	assert.Equal(t, expected.Short, actual.Short)
	assert.Len(t, actual.Commands(), 2)
}
func TestNewCommand_NilLogger(t *testing.T) {
	var log *slog.Logger
	var client app.AuthClient

	cmd := NewCommand(log, client)

	if cmd.Use != "auth" {
		t.Errorf("expected command Use to be 'auth', got %s", cmd.Use)
	}

	if cmd.Short != "Manage user registration, authentication and authorization" {
		t.Errorf("expected command Short to be 'Manage user registration, authentication and authorization', got %s", cmd.Short)
	}

	if len(cmd.Commands()) != 2 {
		t.Errorf("expected 2 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestNewCommand_NilAuthClient(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var client app.AuthClient = nil

	cmd := NewCommand(log, client)

	if cmd == nil {
		t.Error("Expected command to be created, but got nil")
	}

	if cmd.Use != "auth" {
		t.Errorf("Expected command Use to be 'auth', but got %s", cmd.Use)
	}

	if cmd.Short != "Manage user registration, authentication and authorization" {
		t.Errorf("Expected command Short to be 'Manage user registration, authentication and authorization', but got %s", cmd.Short)
	}

	if len(cmd.Commands()) != 2 {
		t.Errorf("Expected command to have 2 subcommands, but got %d", len(cmd.Commands()))
	}
}
