package bin

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

func TestKeepGetBinRunE_NoSecretName(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{
		Env: "dev",
	})

	cmd := NewCommand()
	err := cmd.Flags().Set("name", "")
	require.NoError(t, err)
	err = cmd.Flags().Set("path", "/tmp")
	require.NoError(t, err)

	err = keepGetBinRunE(cmd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "secret name is required")
}
func TestKeepGetBinRunE_SecretPathNotProvided(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{
		Env: "dev",
	})

	cmd := NewCommand()
	err := cmd.Flags().Set("name", "test-secret")
	require.NoError(t, err)
	err = cmd.Flags().Set("path", "")
	require.NoError(t, err)

	err = keepGetBinRunE(cmd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "secret path is required")
}

func TestKeepGetBinRunE_NonExistentSecretName(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{
		Env: "dev",
	})

	cmd := NewCommand()
	err := cmd.Flags().Set("name", "non-existent-secret")
	require.NoError(t, err)
	err = cmd.Flags().Set("path", "/tmp")
	require.NoError(t, err)

	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("GetFile", mock.Anything, "non-existent-secret").Return(nil, fmt.Errorf("secret not found"))

	initClient(mockClient)

	err = keepGetBinRunE(cmd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "secret not found")
	mockClient.AssertExpectations(t)
}
