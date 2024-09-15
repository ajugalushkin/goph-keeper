package bin

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
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

	cmd := &cobra.Command{}
	cmd.Flags().Set("name", "")
	cmd.Flags().Set("path", "/tmp")

	err := keepGetBinRunE(cmd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error reading secret name")
}
func TestKeepGetBinRunE_SecretPathNotProvided(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{
		Env: "dev",
	})

	cmd := NewCommand()
	cmd.Flags().Set("name", "test-secret")
	cmd.Flags().Set("path", "")

	err := keepGetBinRunE(cmd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "secret path is required")
}

func TestKeepGetBinRunE_NonExistentSecretName(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{
		Env: "dev",
	})

	cmd := NewCommand()
	cmd.Flags().Set("name", "non-existent-secret")
	cmd.Flags().Set("path", "/tmp")

	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("GetFile", mock.Anything, "non-existent-secret").Return(nil, fmt.Errorf("secret not found"))

	initClient(mockClient)

	err := keepGetBinRunE(cmd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "secret not found")
	mockClient.AssertExpectations(t)
}

// Successfully reads secret name and path from command flags
func TestKeepGetBinRunE_ReadsSecretNameAndPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "test-secret", "secret name")
	cmd.Flags().String("path", "/tmp", "secret path")

	err := keepGetBinRunE(cmd, []string{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
