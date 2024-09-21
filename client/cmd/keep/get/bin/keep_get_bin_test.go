package bin

import (
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/ajugalushkin/goph-keeper/mocks"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

func TestKeepGetBinRunE_NoSecretName(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{
		Env: "dev",
	})

	cmd := NewCommand()
	getBinCmdFlags(cmd)
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
	getBinCmdFlags(cmd)
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
	getBinCmdFlags(cmd)
	err := cmd.Flags().Set("name", "non-existent-secret")
	require.NoError(t, err)
	err = cmd.Flags().Set("path", "/tmp")
	require.NoError(t, err)

	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("GetFile", mock.Anything, "non-existent-secret", "/tmp").Return(fmt.Errorf("secret not found"))

	initClient(mockClient)

	err = keepGetBinRunE(cmd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "secret not found")
	mockClient.AssertExpectations(t)
}

func TestKeepGetBinRunE_UninitializedClient(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{
		Env: "dev",
	})

	cmd := NewCommand()
	getBinCmdFlags(cmd)
	err := cmd.Flags().Set("name", "test-secret")
	require.NoError(t, err)
	err = cmd.Flags().Set("path", "/tmp/test-secret.txt")
	require.NoError(t, err)

	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("GetFile", mock.Anything, "test-secret", "/tmp/test-secret.txt").Return(nil)

	initClient(mockClient)

	err = keepGetBinRunE(cmd, nil)
	require.NoError(t, err)
	require.Contains(t, "file downloaded: /tmp/test-secret.txt", fmt.Sprintf("%v", err))
	mockClient.AssertExpectations(t)
}

func TestKeepGetBinRunE_TokenCacheLoadFailure(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{
		Env: "dev",
	})

	cmd := NewCommand()
	getBinCmdFlags(cmd)
	err := cmd.Flags().Set("name", "test-secret")
	require.NoError(t, err)
	err = cmd.Flags().Set("path", "/tmp/test-secret.txt")
	require.NoError(t, err)

	token_cache.InitTokenStorage("test/token.txt")

	err = keepGetBinRunE(cmd, nil)
	require.Error(t, err)
}

func TestKeepGetBinRunE_KeeperClientConnectionFailure(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{
		Env: "dev",
	})

	cmd := NewCommand()
	getBinCmdFlags(cmd)
	err := cmd.Flags().Set("name", "test-secret")
	require.NoError(t, err)
	err = cmd.Flags().Set("path", "/tmp/test-secret.txt")
	require.NoError(t, err)

	token_cache.InitTokenStorage("test/token.txt")

	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("GetFile", mock.Anything, "test-secret", "/tmp/test-secret.txt").Return(fmt.Errorf("failed to connect to Keeper service"))

	initClient(mockClient)

	err = keepGetBinRunE(cmd, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to connect to Keeper service")
	mockClient.AssertExpectations(t)
}

func TestKeepGetBinRunE_TokenCacheLoadFailure2(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{
		Env: "dev",
	})

	cmd := NewCommand()
	getBinCmdFlags(cmd)
	err := cmd.Flags().Set("name", "test-secret")
	require.NoError(t, err)
	err = cmd.Flags().Set("path", "/tmp/test-secret.txt")
	require.NoError(t, err)

	token_cache.InitTokenStorage("test_token.txt")
	err = token_cache.GetToken().Save("non-existent-token")
	require.NoError(t, err)

	err = keepGetBinRunE(cmd, nil)
	require.Error(t, err)
}
