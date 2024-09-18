package register

import (
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// Email flag is missing or empty
func TestRegisterCmdRun_MissingEmail(t *testing.T) {
	// Setup
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().StringP("email", "e", "", "User Email")
	cmd.Flags().StringP("password", "p", "password123", "User password")

	// Execute
	err := registerCmdRun(cmd, []string{})

	// Verify
	assert.Error(t, err)
	assert.Equal(t, "email is required", err.Error())
}

func TestRegisterCmdRun_EmptyPassword(t *testing.T) {
	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	registerCmdFlags(cmd)
	err := cmd.Flags().Set("email", "test@example.com")
	require.NoError(t, err)

	err = cmd.Flags().Set("password", "")
	require.NoError(t, err)

	// Act
	err = registerCmdRun(cmd, nil)

	// Assert
	require.Error(t, err)
	require.EqualError(t, err, "password is required")
}

// Successfully retrieves email and password from command-line flags
func TestRegisterCmdRun_Success(t *testing.T) {
	// Setup
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().StringP("email", "e", "test@example.com", "User Email")
	cmd.Flags().StringP("password", "p", "password123", "User password")

	mockClient := new(mocks.AuthClient)
	mockClient.On("Register", mock.Anything, "test@example.com", "password123").Return(nil)

	initClient(mockClient)

	// Execute
	err := registerCmdRun(cmd, []string{})

	// Verify
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRegisterCmdRun_Error(t *testing.T) {
	// Setup
	initClient(nil)
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().StringP("email", "e", "test@example.com", "User Email")
	cmd.Flags().StringP("password", "p", "password123", "User password")

	// Execute
	err := registerCmdRun(cmd, []string{})

	// Verify
	assert.Error(t, err)
}
