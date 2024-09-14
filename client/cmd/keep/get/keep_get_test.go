package get

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

func TestKeepGetRunE_NoSecretName(t *testing.T) {
	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "Secret name")
	errExpected := errors.New("required flag(s) \"name\" not set")

	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("GET")

	// Act
	err := keepGetRunE(cmd, nil)

	// Assert
	assert.EqualError(t, err, errExpected.Error())
}

// Successfully retrieves and decrypts a secret when provided with a valid name
func TestKeepGetRunE_Success(t *testing.T) {
	// Setup
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "valid_secret_name", "secret name")
	args := []string{}

	// Mocking
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	token_cache.InitTokenStorage("test_data/test_token.txt")
	err := token_cache.GetToken().Save("valid_token")
	assert.NoError(t, err)

	mockClient := mocks.NewKeeperClient(t)
	mockClient.On(
		"GetItem",
		mock.Anything,
		&v1.GetItemRequestV1{Name: "valid_secret_name"},
	).Return(&v1.GetItemResponseV1{Content: []byte("encrypted_content")}, nil)
	initClient(mockClient)

	// Execute
	err = keepGetRunE(cmd, args)

	// Assert
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// Handles errors when reading the secret name from command flags
func TestKeepGetRunE_ErrorReadingSecretName(t *testing.T) {
	// Setup
	cmd := &cobra.Command{}
	args := []string{}

	// Mocking
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd.Flags().String("name", "", "secret name")

	// Execute
	err := keepGetRunE(cmd, args)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error reading secret name")
}
