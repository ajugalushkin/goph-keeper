package text

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
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/ajugalushkin/goph-keeper/client/secret"
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

func TestKeeperUpdateTextCmdRunE_EmptySecretName(t *testing.T) {
	initClient(nil)
	token_cache.InitTokenStorage("")

	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	// Set up the context and logger
	name := ""
	data := "test-data"

	// Mock dependencies
	token_cache.InitTokenStorage("test_token.txt")
	err := token_cache.GetToken().Save("test-token")
	assert.NoError(t, err)

	config.InitConfig(&config.Config{Client: config.Client{Address: "test-address"}})

	// Create a command and set up flags
	cmd := NewCommand()
	cmd.SetArgs([]string{"--name", name, "--data", data})

	// Act
	err = keeperUpdateTextCmdRunE(cmd, nil)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "secret name cannot be empty")
}

// Handles error when the secret name flag is missing or empty
func TestKeeperUpdateTextCmdRunE_MissingName(t *testing.T) {
	initClient(nil)
	// Setup
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "secret name")
	cmd.Flags().String("data", "test-data", "secret data")

	// Execute
	err := keeperUpdateTextCmdRunE(cmd, []string{})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "secret name cannot be empty", err.Error())
}

func TestKeeperUpdateTextCmdRunE_NonExistentSecretName(t *testing.T) {
	initClient(nil)
	token_cache.InitTokenStorage("")

	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	// Set up the context and logger
	name := "non-existent-secret"
	data := "test-data"

	// Mock dependencies
	token_cache.InitTokenStorage("token.txt")
	err := token_cache.GetToken().Save("test-token")
	assert.NoError(t, err)

	config.InitConfig(&config.Config{Client: config.Client{Address: "test-address"}})

	mockClient := mocks.NewKeeperClient(t)

	// Create a Text secret object with the provided data
	text := vaulttypes.Text{
		Data: data,
	}

	// Encrypt the secret data
	content, err := secret.NewCryptographer().Encrypt(text)
	assert.NoError(t, err)

	mockClient.On(
		"UpdateItem",
		mock.Anything,
		&keeperv1.UpdateItemRequestV1{Name: name, Content: content},
	).Return(nil, fmt.Errorf("failed to update secret"))
	initClient(mockClient)

	// Create a command and set up flags
	cmd := &cobra.Command{}
	cmd.Flags().String("name", name, "secret name")
	cmd.Flags().String("data", data, "secret data")

	// Act
	err = keeperUpdateTextCmdRunE(cmd, nil)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update secret")
}

func TestKeeperUpdateTextCmdRunE_Success(t *testing.T) {
	initClient(nil)
	initCipher(nil)
	token_cache.InitTokenStorage("")

	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	// Set up the context and logger
	name := "non-existent-secret"
	data := "test-data"

	// Mock dependencies
	token_cache.InitTokenStorage("token.txt")
	err := token_cache.GetToken().Save("test-token")
	assert.NoError(t, err)

	config.InitConfig(&config.Config{Client: config.Client{Address: "test-address"}})

	mockClient := mocks.NewKeeperClient(t)

	// Create a Text secret object with the provided data
	text := vaulttypes.Text{
		Data: data,
	}

	// Encrypt the secret data
	content, err := secret.NewCryptographer().Encrypt(text)
	assert.NoError(t, err)

	mockClient.On(
		"UpdateItem",
		mock.Anything,
		&keeperv1.UpdateItemRequestV1{Name: name, Content: content},
	).Return(&keeperv1.UpdateItemResponseV1{Name: name, Version: "1"}, nil)
	initClient(mockClient)

	// Create a command and set up flags
	cmd := &cobra.Command{}
	cmd.Flags().String("name", name, "secret name")
	cmd.Flags().String("data", data, "secret data")

	// Act
	err = keeperUpdateTextCmdRunE(cmd, nil)

	// Assert
	assert.NoError(t, err)
}

func TestKeeperUpdateTextCmdRunE_Error(t *testing.T) {
	initClient(nil)
	initCipher(nil)
	token_cache.InitTokenStorage("")

	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	// Set up the context and logger
	name := "non-existent-secret"
	data := "test-data"

	// Mock dependencies
	token_cache.InitTokenStorage("token.txt")
	err := token_cache.GetToken().Save("test-token")
	assert.NoError(t, err)

	config.InitConfig(&config.Config{Client: config.Client{Address: "test-address"}})

	// Create a command and set up flags
	cmd := &cobra.Command{}
	cmd.Flags().String("name", name, "secret name")
	cmd.Flags().String("data", data, "secret data")

	// Act
	err = keeperUpdateTextCmdRunE(cmd, nil)

	// Assert
	assert.Error(t, err)
}

func TestKeeperUpdateTextCmdRunE_Error2(t *testing.T) {
	initClient(nil)
	initCipher(nil)
	token_cache.InitTokenStorage("")

	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	// Set up the context and logger
	name := "non-existent-secret"
	data := "test-data"

	// Mock dependencies
	token_cache.InitTokenStorage("token.txt")
	err := token_cache.GetToken().Save("test-token")
	assert.NoError(t, err)

	config.InitConfig(&config.Config{Client: config.Client{Address: "test-address"}})

	mockCipher := mocks.NewCipher(t)
	mockCipher.On("Encrypt", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to encrypt data"))
	initCipher(mockCipher)

	// Create a command and set up flags
	cmd := &cobra.Command{}
	cmd.Flags().String("name", name, "secret name")
	cmd.Flags().String("data", data, "secret data")

	// Act
	err = keeperUpdateTextCmdRunE(cmd, nil)

	// Assert
	assert.Error(t, err)
}
