package creds

import (
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// Fails to retrieve 'name' flag from command-line arguments
func TestFailToRetrieveNameFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("login", "testLogin", "login flag")
	cmd.Flags().String("password", "testPassword", "password flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

// Fails to retrieve 'login' flag from command-line arguments
func TestFailToRetrieveLoginFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name flag")
	cmd.Flags().String("password", "testPassword", "password flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected an error for failing to retrieve 'login' flag, but no error was returned")
	}
}

// Fails to retrieve 'password' flag from command-line arguments
func TestFailToRetrievePasswordFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name flag")
	cmd.Flags().String("login", "testLogin", "login flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error for failing to retrieve 'password' flag, but got no error")
	}
}

func TestCreateCredentialsWithEmptyNameFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "name flag")
	cmd.Flags().String("login", "testLogin", "login flag")
	cmd.Flags().String("password", "testPassword", "password flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error for empty 'name' flag, but no error was returned")
	}
}
func TestEmptyLoginFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name flag")
	cmd.Flags().String("login", "", "login flag")
	cmd.Flags().String("password", "testPassword", "password flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error for empty 'login' flag, but no error was returned")
	}
}
func TestEmptyPasswordFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name flag")
	cmd.Flags().String("login", "testLogin", "login flag")
	cmd.Flags().String("password", "", "password flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error for empty 'password' flag, but got nil")
	}
}
