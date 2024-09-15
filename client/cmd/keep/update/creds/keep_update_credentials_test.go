package creds

import (
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
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
