package text

import (
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// Successfully reads secret name and data from command flags
func TestKeepCreateTextCmdRunE_Success(t *testing.T) {
	// Setup
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "test-secret", "secret name")
	cmd.Flags().String("data", "test-data", "secret data")

	// Mock logger
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Mock KeeperClient
	mockClient := mocks.NewKeeperClient(t)
	client = mockClient

	// Mock response
	mockResp := &v1.CreateItemResponseV1{
		Name:    "test-secret",
		Version: "1",
	}
	mockClient.On("CreateItem", mock.Anything, mock.Anything).Return(mockResp, nil)

	// Execute
	err := keepCreateTextCmdRunE(cmd, []string{})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "test-secret", mockResp.GetName())
	assert.Equal(t, "1", mockResp.GetVersion())
}

// Command flags for name or data are missing or invalid
func TestKeepCreateTextCmdRunE_MissingFlags(t *testing.T) {
	// Setup
	cmd := &cobra.Command{}

	// Mock logger
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Execute with missing name flag
	err := keepCreateTextCmdRunE(cmd, []string{})

	// Assert
	assert.Error(t, err)

	// Add name flag but missing data flag
	cmd.Flags().String("name", "test-secret", "secret name")

	// Execute with missing data flag
	err = keepCreateTextCmdRunE(cmd, []string{})

	// Assert
	assert.Error(t, err)
}
