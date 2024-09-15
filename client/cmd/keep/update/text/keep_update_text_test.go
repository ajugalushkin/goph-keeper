package text

import (
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
)

func TestKeeperUpdateTextCmdRunE_EmptySecretName(t *testing.T) {
	// Arrange
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	// Set up the context and logger
	name := ""
	data := "test-data"

	// Mock dependencies
	token_cache.InitTokenStorage("test_data/token.txt")
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
