package auth

import (
	"log/slog"
	"os"
	"testing"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
)

func TestNewCommand_ValidConfig(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := &config.Config{}

	cmd := NewCommand(log, cfg)

	if cmd == nil {
		t.Fatal("Expected NewCommand to return a non-nil command")
	}

	if cmd.Use != "auth" {
		t.Errorf("Expected command Use to be 'auth', got %s", cmd.Use)
	}

	if cmd.Short != "Manage user registration, authentication and authorization" {
		t.Errorf("Expected command Short to be 'Manage user registration, authentication and authorization', got %s", cmd.Short)
	}

	if len(cmd.Commands()) != 2 {
		t.Fatalf("Expected command to have 2 subcommands, got %d", len(cmd.Commands()))
	}

	loginCmd := cmd.Commands()[0]
	if loginCmd.Name() != "login" {
		t.Errorf("Expected first subcommand to be 'login', got %s", loginCmd.Name())
	}

	registerCmd := cmd.Commands()[1]
	if registerCmd.Name() != "register" {
		t.Errorf("Expected second subcommand to be 'register', got %s", registerCmd.Name())
	}
}
