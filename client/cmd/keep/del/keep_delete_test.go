package del

import (
	"errors"
	"github.com/ajugalushkin/goph-keeper/mocks"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// Handles error when secret name flag is missing or invalid
func TestKeepDeleteCmdRunE_MissingNameFlag(t *testing.T) {
	// Setup
	cmd := &cobra.Command{}
	args := []string{}

	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Execute
	err := keepDeleteCmdRunE(cmd, args)

	// Verify
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading secret name")
}

func TestKeepDeleteCmdRunE_NonExistentSecret(t *testing.T) {
	// Arrange
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "non-existent-secret", "secret name")
	args := []string{}

	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	token_cache.InitTokenStorage("test/test_token.txt")
	err := token_cache.GetToken().Save("mock-token")
	require.NoError(t, err)

	mockKeeperClient := mocks.NewKeeperClient(t)
	mockKeeperClient.On(
		"DeleteItem",
		mock.Anything,
		&keeperv1.DeleteItemRequestV1{Name: "non-existent-secret"},
	).Return(nil, errors.New("secret not found"))

	initClient(mockKeeperClient)

	// Execute
	err = keepDeleteCmdRunE(cmd, args)

	// Verify
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "secret not found")
}
func TestKeepDeleteCmdRunE_EmptyTokenCache(t *testing.T) {
	// Arrange
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "test-secret", "secret name")
	args := []string{}

	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	token_cache.InitTokenStorage("test/empty_token.txt")

	mockKeeperClient := mocks.NewKeeperClient(t)
	mockKeeperClient.On(
		"DeleteItem",
		mock.Anything,
		&keeperv1.DeleteItemRequestV1{Name: "test-secret"},
	).Return(nil, errors.New("empty token cache"))

	initClient(mockKeeperClient)

	// Execute
	err := keepDeleteCmdRunE(cmd, args)

	// Verify
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty token cache")
}
func TestKeepDeleteCmdRunE_EmptyNameFlag(t *testing.T) {
	// Arrange
	cmd := &cobra.Command{}
	args := []string{}

	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Execute
	err := keepDeleteCmdRunE(cmd, args)

	// Verify
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading secret name")
}

func TestKeepDeleteCmdRunE_Success(t *testing.T) {
	// Arrange
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "test-secret", "secret name")
	args := []string{}

	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	token_cache.InitTokenStorage("test/empty_token.txt")

	mockKeeperClient := mocks.NewKeeperClient(t)
	mockKeeperClient.On(
		"DeleteItem",
		mock.Anything,
		&keeperv1.DeleteItemRequestV1{Name: "test-secret"},
	).Return(nil, nil)

	initClient(mockKeeperClient)

	// Execute
	err := keepDeleteCmdRunE(cmd, args)

	// Verify
	assert.NoError(t, err)
}
