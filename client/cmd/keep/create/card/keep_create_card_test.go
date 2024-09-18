package card

import (
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// Fails to read the secret name flag
func TestCreateCardCmdRunE_FailsToReadSecretNameFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := NewCommand()
	cardCmdFlags(cmd)

	err := cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	err = createCardCmdRunE(cmd, []string{})
	assert.Error(t, err, "name is required")
}

// Fails to read the card number flag
func TestCreateCardCmdRunE_FailsToReadCardNumberFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := NewCommand()
	cardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	err = createCardCmdRunE(cmd, []string{})
	assert.Error(t, err, "name is required")

	err = createCardCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatal("card number is required")
	}
}

// Fails to read the card expiry date flag
func TestCreateCardCmdRunE_FailsToReadExpiryDateFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := NewCommand()
	cardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	err = createCardCmdRunE(cmd, []string{})
	assert.Error(t, err, "expiry date is required")
}

func TestCreateCardCmdRunE_FailsToReadCodeFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := NewCommand()
	cardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("holder", "John Doe")
	require.NoError(t, err)

	err = createCardCmdRunE(cmd, []string{})
	assert.Error(t, err, "security code is required")
}

func TestCreateCardCmdRunE_FailsToReadHolderFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := NewCommand()
	cardCmdFlags(cmd)

	err := cmd.Flags().Set("name", "test_secret")
	require.NoError(t, err)

	err = cmd.Flags().Set("number", "1234567890123456")
	require.NoError(t, err)

	err = cmd.Flags().Set("date", "12/24")
	require.NoError(t, err)

	err = cmd.Flags().Set("code", "123")
	require.NoError(t, err)

	err = createCardCmdRunE(cmd, []string{})
	assert.Error(t, err, "card holder is required")
}

func TestCreateCardCmdRunE_SucceedsWhenAllFlagsAreProvided(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := NewCommand()
	cardCmdFlags(cmd)

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

	token_cache.InitTokenStorage("./test/token.txt")

	err = createCardCmdRunE(cmd, []string{})
	assert.Error(t, err)
}
