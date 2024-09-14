package auth

import (
	"log/slog"
	"os"
	"testing"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
)

func TestNewCommand_NilConfig(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var cfg config.Client

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
		t.Errorf("Expected command to have 2 subcommands, got %d", len(cmd.Commands()))
	}
}
