package card

import (
	"context"
	"errors"
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/mocks"
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

func TestKeeperUpdateCardCmdRunE_MissingNameFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	// Set up the mock Keeper client
	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	// Create a Cobra command and set the flags
	cmd := &cobra.Command{}
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

	cmd := &cobra.Command{}
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
	cmd := &cobra.Command{}
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
	cmd := &cobra.Command{}
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
	cmd := &cobra.Command{}
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
	initClient(nil)
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
	initClient(nil)
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

func TestKeeperUpdateCardCmdRunE_Success(t *testing.T) {
	initClient(nil)
	initCipher(nil)

	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	newMockCipher := mocks.NewCipher(t)

	card := vaulttypes.Card{
		Number:       "1234567890123456",
		ExpiryDate:   "12/24",
		SecurityCode: "123",
		Holder:       "John Doe",
	}

	content := []byte("encrypted_data")
	newMockCipher.On("Encrypt", card).Return(content, nil)
	initCipher(newMockCipher)

	// Set up the mock Keeper client
	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("UpdateItem", context.Background(), &v1.UpdateItemRequestV1{
		Name:    "test_secret",
		Content: content,
	}).Return(&v1.UpdateItemResponseV1{
		Name:    "test_secret",
		Version: "1",
	}, nil)
	initClient(mockClient)

	// Create a Cobra command and set the flags
	cmd := &cobra.Command{}
	updateCardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	// Run the command with empty security code
	err = keeperUpdateCardCmdRunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestKeeperUpdateCardCmdRunE_Error(t *testing.T) {
	initClient(nil)
	initCipher(nil)

	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	newMockCipher := mocks.NewCipher(t)

	card := vaulttypes.Card{
		Number:       "1234567890123456",
		ExpiryDate:   "12/24",
		SecurityCode: "123",
		Holder:       "John Doe",
	}

	content := []byte("encrypted_data")
	newMockCipher.On("Encrypt", card).Return(content, nil)
	initCipher(newMockCipher)

	// Set up the mock Keeper client
	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("UpdateItem", context.Background(), &v1.UpdateItemRequestV1{
		Name:    "test_secret",
		Content: content,
	}).Return(nil, errors.New("error updating card"))
	initClient(mockClient)

	// Create a Cobra command and set the flags
	cmd := &cobra.Command{}
	updateCardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	// Run the command with empty security code
	err = keeperUpdateCardCmdRunE(cmd, []string{})
	assert.Error(t, err)
}

func TestKeeperUpdateCardCmdRunE_Error2(t *testing.T) {
	initClient(nil)
	initCipher(nil)
	
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	newMockCipher := mocks.NewCipher(t)

	newMockCipher.On("Encrypt", mock.Anything).Return(nil, errors.New("error encrypting card"))
	initCipher(newMockCipher)

	// Create a Cobra command and set the flags
	cmd := &cobra.Command{}
	updateCardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	// Run the command with empty security code
	err = keeperUpdateCardCmdRunE(cmd, []string{})
	assert.Error(t, err)
}
