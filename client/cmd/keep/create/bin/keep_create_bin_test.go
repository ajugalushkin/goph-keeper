package bin

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

func TestKeepCreateBinCmdRunE_EmptySecretName(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	cmd := NewCommand()
	err := cmd.Flags().Set("name", "")
	require.NoError(t, err)

	err = cmd.Flags().Set("file_path", "/path/to/file")
	require.NoError(t, err)

	err = keepCreateBinCmdRunE(cmd, nil)
	require.Error(t, err)
	require.EqualError(t, err, "name is required")
}
func TestKeepCreateBinCmdRunE_EmptyFilePath(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	cmd := NewCommand()

	err := cmd.Flags().Set("name", "test-secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("file_path", "")
	require.NoError(t, err)

	ctx := context.Background()
	clientMock := &mocks.KeeperClient{}
	initClient(clientMock)

	clientMock.On(
		"CreateItemStream",
		ctx,
		"test-secret",
		"",
	).Return(nil, fmt.Errorf("file path is required")).Maybe()

	err = keepCreateBinCmdRunE(cmd, nil)
	require.Error(t, err)
	require.EqualError(t, err, "file_path is required")

	clientMock.AssertExpectations(t)
}
func TestKeepCreateBinCmdRunE_NonExistentFilePath(t *testing.T) {
	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	clientMock := mocks.NewKeeperClient(t)
	clientMock.On(
		"CreateItemStream",
		context.Background(),
		"test-secret",
		"/non/existent/path",
	).Return(nil, fmt.Errorf("open /non/existent/path: no such file or directory")).Maybe()
	initClient(clientMock)

	cmd := NewCommand()

	err := cmd.Flags().Set("name", "test-secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("file_path", "/non/existent/path")
	require.NoError(t, err)

	expectedErr := fmt.Errorf("open /non/existent/path: no such file or directory")

	// Act
	err = keepCreateBinCmdRunE(cmd, nil)

	// Assert
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error: %v, got: %v", expectedErr, err)
	}
}
func TestKeepCreateBinCmdRunE_InvalidToken(t *testing.T) {
	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	ctx := context.Background()
	name := "test-secret"
	filePath := "test-file.bin"
	invalidToken := "invalid-token"
	expectedError := "invalid token"

	// Mock the keeper client
	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("CreateItemStream", ctx, name, filePath).Return(nil, fmt.Errorf(expectedError))

	// Initialize the command with the mock client
	initClient(mockClient)

	// Create a new command and set the flags
	cmd := NewCommand()
	cmd.SetArgs([]string{"--name", name, "--file_path", filePath})
	err := cmd.Flags().Set("token", invalidToken)
	require.NoError(t, err)

	// Act
	err = cmd.ExecuteContext(ctx)

	// Assert
	assert.ErrorContains(t, err, expectedError)
	mockClient.AssertExpectations(t)
}
