package card

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

func TestKeeperUpdateCardCmdRunE_MissingNameFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Set up the mock Keeper client
	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	// Create a Cobra command and set the flags
	cmd := NewCommand()
	updateCardCmdFlags(cmd)

	err := cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	// Run the command with missing name flag
	err = keeperUpdateCardCmdRunE(cmd, []string{})

	// Verify the expected error
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid secret name")

	// Verify that the mock client was not called
	mockClient.AssertNotCalled(t, "UpdateItem", mock.Anything, mock.Anything)
}
func TestKeeperUpdateCardCmdRunE_MissingNumberFlag(t *testing.T) {
	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := NewCommand()
	updateCardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	// Act
	err = keeperUpdateCardCmdRunE(cmd, nil)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid card number")
}

func TestKeeperUpdateCardCmdRunE_MissingExpiryDateFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Set up the mock Keeper client
	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	// Create a Cobra command and set the flags
	cmd := NewCommand()
	updateCardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	// Run the command with missing expiry date flag
	err = keeperUpdateCardCmdRunE(cmd, []string{})

	// Verify the expected error
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid card expiry date")

	// Verify that the mock client was not called
	mockClient.AssertNotCalled(t, "UpdateItem", mock.Anything, mock.Anything)
}

func TestKeeperUpdateCardCmdRunE_MissingSecurityCodeFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Set up the mock Keeper client
	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	// Create a Cobra command and set the flags
	cmd := NewCommand()
	updateCardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	// Run the command with missing security code flag
	err = keeperUpdateCardCmdRunE(cmd, []string{})

	// Verify the expected error
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid card security code")

	// Verify that the mock client was not called
	mockClient.AssertNotCalled(t, "UpdateItem", mock.Anything, mock.Anything)
}

func TestKeeperUpdateCardCmdRunE_MissingHolderFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Set up the mock Keeper client
	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	// Create a Cobra command and set the flags
	cmd := NewCommand()
	updateCardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	// Run the command with missing holder flag
	err = keeperUpdateCardCmdRunE(cmd, []string{})

	// Verify the expected error
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid card holder")

	// Verify that the mock client was not called
	mockClient.AssertNotCalled(t, "UpdateItem", mock.Anything, mock.Anything)
}

func TestKeeperUpdateCardCmdRunE_EmptyCardNumber(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Set up the mock Keeper client
	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	// Create a Cobra command and set the flags
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "test-card", "")
	cmd.Flags().String("number", "", "")
	cmd.Flags().String("date", "12/24", "")
	cmd.Flags().String("code", "123", "")
	cmd.Flags().String("holder", "John Doe", "")

	// Run the command with empty card number
	err := keeperUpdateCardCmdRunE(cmd, []string{})

	// Verify the expected error
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid card number")

	// Verify that the mock client was not called
	mockClient.AssertNotCalled(t, "UpdateItem", mock.Anything, mock.Anything)
}
func TestKeeperUpdateCardCmdRunE_EmptyExpiryDate(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Set up the mock Keeper client
	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	// Create a Cobra command and set the flags
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "test-card", "")
	cmd.Flags().String("number", "1234567890123456", "")
	cmd.Flags().String("date", "", "")
	cmd.Flags().String("code", "123", "")
	cmd.Flags().String("holder", "John Doe", "")

	// Run the command with empty expiry date
	err := keeperUpdateCardCmdRunE(cmd, []string{})

	// Verify the expected error
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid card expiry date")

	// Verify that the mock client was not called
	mockClient.AssertNotCalled(t, "UpdateItem", mock.Anything, mock.Anything)
}
func TestKeeperUpdateCardCmdRunE_EmptySecurityCode(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Set up the mock Keeper client
	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	// Create a Cobra command and set the flags
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "test-secret", "")
	cmd.Flags().String("number", "1234567890123456", "")
	cmd.Flags().String("date", "12/24", "")
	cmd.Flags().String("code", "", "")
	cmd.Flags().String("holder", "John Doe", "")

	// Run the command with empty security code
	err := keeperUpdateCardCmdRunE(cmd, []string{})

	// Verify the expected error
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid card security code")

	// Verify that the mock client was not called
	mockClient.AssertNotCalled(t, "UpdateItem", mock.Anything, mock.Anything)
}
