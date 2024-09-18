package login

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
)

func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()
	return strings.TrimSpace(buf.String()), err
}

func TestAuthLoginCmdRunE_ValidEmailAndPassword(t *testing.T) {
	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	email := "test@example.com"
	password := "testpassword"
	expectedToken := "test_access_token"

	mockAuthClient := mocks.NewAuthClient(t)
	mockAuthClient.On("Login", context.Background(), email, password).Return(expectedToken, nil).Maybe()
	initClient(mockAuthClient)

	token_cache.InitTokenStorage("./test_data/test_save_token.txt")
	assert.NoError(t, token_cache.GetToken().Save(expectedToken))

	cmd := NewCommand()
	loginCmdFlags(cmd)

	err := cmd.Flags().Set("email", email)
	require.NoError(t, err)

	err = cmd.Flags().Set("password", password)
	require.NoError(t, err)

	// Act
	err = authLoginCmdRunE(cmd, nil)

	// Assert
	assert.Nil(t, err)
	mockAuthClient.AssertExpectations(t)
	fmt.Printf("Access Token: %s\n", expectedToken)
}

func TestAuthLoginCmdRunE_EmptyEmail(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	cmd := NewCommand()
	loginCmdFlags(cmd)
	err := cmd.Flags().Set("email", "")
	require.NoError(t, err)

	err = cmd.Flags().Set("password", "testpassword")
	require.NoError(t, err)

	err = authLoginCmdRunE(cmd, nil)
	require.Error(t, err)

	loginCmd := NewCommand()

	loginCmd.SetArgs([]string{"password", "testpassword"})

	err = authLoginCmdRunE(loginCmd, nil)
	require.Error(t, err)
}

func TestAuthLoginCmdRunE_EmptyPassword(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	cmd := NewCommand()
	loginCmdFlags(cmd)

	email := "test@example.com"
	password := ""

	err := cmd.Flags().Set("email", email)
	require.NoError(t, err)
	err = cmd.Flags().Set("password", password)
	require.NoError(t, err)

	err = authLoginCmdRunE(cmd, nil)
	require.Error(t, err)
	require.EqualError(t, err, "password is required")

	loginCmd := NewCommand()
	loginCmdFlags(loginCmd)

	err = loginCmd.Flags().Set("email", email)
	require.NoError(t, err)

	err = authLoginCmdRunE(loginCmd, nil)
	require.Error(t, err)
	require.EqualError(t, err, "password is required")
}
