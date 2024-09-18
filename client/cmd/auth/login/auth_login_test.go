package login

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
)

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

	cmd := &cobra.Command{}
	loginCmdFlags(cmd)

	err := cmd.Flags().Set("email", "")
	require.NoError(t, err)

	err = cmd.Flags().Set("password", "testpassword")
	require.NoError(t, err)

	err = authLoginCmdRunE(cmd, nil)
	require.Error(t, err)
	require.EqualError(t, err, "email is required")
}
func TestAuthLoginCmdRunE_EmptyPassword(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	cmd := &cobra.Command{}
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
}
