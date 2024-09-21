package get

import (
	"fmt"
	"github.com/ajugalushkin/goph-keeper/mocks"
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	secretCipher "github.com/ajugalushkin/goph-keeper/client/secret/mocks"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// TestKeepGetRunE_SecretNotFound verifies behavior when the secret name does not exist in the vault
func TestKeepGetRunE_SecretNotFound(t *testing.T) {
	// Setup
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "non-existent-secret", "secret name")
	args := []string{}

	mockClient := mocks.NewKeeperClient(t)
	client = mockClient

	mockClient.On(
		"GetItem",
		mock.Anything,
		&v1.GetItemRequestV1{Name: "non-existent-secret"},
	).Return(nil, fmt.Errorf("secret not found"))

	// Execute
	err := keepGetRunE(cmd, args)

	// Verify
	assert.Error(t, err)
	assert.EqualError(t, err, "secret not found")
	mockClient.AssertExpectations(t)
}

// Handles errors when reading the secret name from command flags
func TestKeepGetRunE_ErrorReadingSecretName(t *testing.T) {
	// Setup
	cmd := &cobra.Command{}
	args := []string{}

	// Simulate error in reading the secret name
	cmd.Flags().String("name", "", "secret name")

	// Execute
	err := keepGetRunE(cmd, args)

	// Verify
	assert.Error(t, err)
	assert.Equal(t, "secret name is required", err.Error())
}
func TestKeepGetRunE_Error(t *testing.T) {
	// Setup
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "non-existent-secret", "secret name")
	args := []string{}

	mockClient := mocks.NewKeeperClient(t)
	client = mockClient

	mockClient.On(
		"GetItem",
		mock.Anything,
		&v1.GetItemRequestV1{Name: "non-existent-secret"},
	).Return(nil, nil)

	mockCipher := secretCipher.NewCipher(t)
	cipher = mockCipher

	mockCipher.On("Decrypt", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("decryption error"))

	// Execute
	err := keepGetRunE(cmd, args)

	// Verify
	assert.Error(t, err)
	assert.EqualError(t, err, "decryption error")
	mockCipher.AssertExpectations(t)
}
