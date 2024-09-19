package creds

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/secret"
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

func TestKeepUpdateCredsRunE_SecretNameNotProvided(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	updateCredsCmdFlags(cmd)

	err := cmd.Flags().Set("login", "test_login")
	require.NoError(t, err)
	err = cmd.Flags().Set("password", "test_password")
	require.NoError(t, err)

	err = keepUpdateCredsRunE(cmd, []string{})
	require.Error(t, err)
	require.EqualError(t, err, "name is required")
}

func TestKeepUpdateCredsRunE_NoLogin(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	updateCredsCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test-secret")
	require.NoError(t, err)
	err = cmd.Flags().Set("password", "test-password")
	require.NoError(t, err)

	err = keepUpdateCredsRunE(cmd, nil)
	require.Error(t, err)
	assert.EqualError(t, err, "login is required")
}
func TestKeepUpdateCredsRunE_PasswordNotProvided(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	updateCredsCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test-secret")
	require.NoError(t, err)
	err = cmd.Flags().Set("login", "test-user")
	require.NoError(t, err)
	err = cmd.Flags().Set("password", "")
	require.NoError(t, err)

	err = keepUpdateCredsRunE(cmd, nil)
	require.Error(t, err)
	assert.EqualError(t, err, "password is required")
}

func TestKeepUpdateCredsRunE_EmptyName(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	updateCredsCmdFlags(cmd)

	err := cmd.Flags().Set("name", "")
	require.NoError(t, err)
	err = cmd.Flags().Set("login", "test_login")
	require.NoError(t, err)
	err = cmd.Flags().Set("password", "test_password")
	require.NoError(t, err)

	err = keepUpdateCredsRunE(cmd, nil)
	require.Error(t, err)
	assert.EqualError(t, err, "name is required")
}

// Successfully reads command-line flags for name, login, and password
func TestKeepUpdateCredsRunE_SuccessfulFlagRead(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name of the secret")
	cmd.Flags().String("login", "testLogin", "login for the secret")
	cmd.Flags().String("password", "testPassword", "password for the secret")

	mockClient := mocks.NewKeeperClient(t)

	credentials := vaulttypes.Credentials{
		Login:    "testLogin",
		Password: "testPassword",
	}

	content, err := secret.NewCryptographer().Encrypt(credentials)
	assert.NoError(t, err)

	mockClient.On("UpdateItem", context.Background(), &v1.UpdateItemRequestV1{
		Name:    "testName",
		Content: content,
	}).Return(&v1.UpdateItemResponseV1{Name: "testName", Version: "1"}, nil)
	initClient(mockClient)

	err = keepUpdateCredsRunE(cmd, []string{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// Command-line flag for name is missing or empty
func TestKeepUpdateCredsRunE_MissingNameFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("login", "testLogin", "login for the secret")
	cmd.Flags().String("password", "testPassword", "password for the secret")

	err := keepUpdateCredsRunE(cmd, []string{})
	if err == nil || err.Error() != "error reading secret name" {
		t.Fatalf("Expected 'name is required' error, got %v", err)
	}
}

// Creates Credentials object with provided login and password
func TestKeepUpdateCredsRunE_CreatesCredentials(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name of the secret")
	cmd.Flags().String("login", "testLogin", "login for the secret")
	cmd.Flags().String("password", "testPassword", "password for the secret")

	err := keepUpdateCredsRunE(cmd, []string{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// Command-line flag for login is missing or empty
func TestKeepUpdateCredsRunE_LoginFlagMissingOrEmpty(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name of the secret")
	cmd.Flags().String("login", "", "login for the secret")
	cmd.Flags().String("password", "testPassword", "password for the secret")

	err := keepUpdateCredsRunE(cmd, []string{})
	if err == nil || err.Error() != "login is required" {
		t.Fatalf("Expected error 'login is required', got %v", err)
	}
}

// Command-line flag for password is missing or empty
func TestKeepUpdateCredsRunE_PasswordFlagMissingOrEmpty(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name of the secret")
	cmd.Flags().String("login", "testLogin", "login for the secret")
	cmd.Flags().String("password", "", "password for the secret")

	err := keepUpdateCredsRunE(cmd, []string{})
	if err == nil || err.Error() != "password is required" {
		t.Fatalf("Expected error 'password is required', got %v", err)
	}
}
