package card

import (
	"context"
	"errors"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/ajugalushkin/goph-keeper/client/secret"
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/mocks"
	"github.com/spf13/cobra"
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

	cmd := &cobra.Command{}
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

	cmd := &cobra.Command{}
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

	cmd := &cobra.Command{}
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

	cmd := &cobra.Command{}
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

	cmd := &cobra.Command{}
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

func TestCreateCardCmdRunE_ErrorWhenAllFlagsAreProvided(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
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

func TestCreateCardCmdRunE_SuccessWhenAllFlagsAreProvided(t *testing.T) {
	initClient(nil)
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cardCmdFlags(cmd)

	name := "test_secret"
	err := cmd.Flags().Set("name", name)
	require.NoError(t, err)

	cardNumber := "1234567890123456"
	err = cmd.Flags().Set("number", cardNumber)
	require.NoError(t, err)

	date := "12/24"
	err = cmd.Flags().Set("date", date)
	require.NoError(t, err)

	code := "123"
	err = cmd.Flags().Set("code", code)
	require.NoError(t, err)

	holder := "John Doe"
	err = cmd.Flags().Set("holder", holder)
	require.NoError(t, err)

	token_cache.InitTokenStorage("./test/token.txt")

	card := vaulttypes.Card{
		Number:       cardNumber,
		ExpiryDate:   date,
		SecurityCode: code,
		Holder:       holder,
	}

	// Encrypt the card details.
	content, err := secret.NewCryptographer().Encrypt(card)
	require.NoError(t, err)

	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("CreateItem", context.Background(), &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	}).Return(&v1.CreateItemResponseV1{
		Name:    name,
		Version: "1",
	}, nil)
	initClient(mockClient)

	err = createCardCmdRunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestCreateCardCmdRunE_Error(t *testing.T) {
	initClient(nil)
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cardCmdFlags(cmd)

	name := "test_secret"
	err := cmd.Flags().Set("name", name)
	require.NoError(t, err)

	cardNumber := "1234567890123456"
	err = cmd.Flags().Set("number", cardNumber)
	require.NoError(t, err)

	date := "12/24"
	err = cmd.Flags().Set("date", date)
	require.NoError(t, err)

	code := "123"
	err = cmd.Flags().Set("code", code)
	require.NoError(t, err)

	holder := "John Doe"
	err = cmd.Flags().Set("holder", holder)
	require.NoError(t, err)

	token_cache.InitTokenStorage("./test/token.txt")

	card := vaulttypes.Card{
		Number:       cardNumber,
		ExpiryDate:   date,
		SecurityCode: code,
		Holder:       holder,
	}

	// Encrypt the card details.
	content, err := secret.NewCryptographer().Encrypt(card)
	require.NoError(t, err)

	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("CreateItem", context.Background(), &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	}).Return(&v1.CreateItemResponseV1{
		Name:    name,
		Version: "1",
	}, errors.New("Internal server error"))
	initClient(mockClient)

	err = createCardCmdRunE(cmd, []string{})
	assert.Error(t, err)
}

func TestCreateCardCmdRunE_Error2(t *testing.T) {
	initClient(nil)
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cardCmdFlags(cmd)

	name := "test_secret"
	err := cmd.Flags().Set("name", name)
	require.NoError(t, err)

	cardNumber := "1234567890123456"
	err = cmd.Flags().Set("number", cardNumber)
	require.NoError(t, err)

	date := "12/24"
	err = cmd.Flags().Set("date", date)
	require.NoError(t, err)

	code := "123"
	err = cmd.Flags().Set("code", code)
	require.NoError(t, err)

	holder := "John Doe"
	err = cmd.Flags().Set("holder", holder)
	require.NoError(t, err)

	token_cache.InitTokenStorage("./test/token.txt")

	card := vaulttypes.Card{
		Number:       cardNumber,
		ExpiryDate:   date,
		SecurityCode: code,
		Holder:       holder,
	}

	// Encrypt the card details.
	content, err := secret.NewCryptographer().Encrypt(card)
	require.NoError(t, err)

	mockClient := mocks.NewKeeperClient(t)
	mockClient.On("CreateItem", context.Background(), &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	}).Return(nil, nil)
	initClient(mockClient)

	err = createCardCmdRunE(cmd, []string{})
	assert.Error(t, err)
}

func TestCreateCardCmdRunE_Error3(t *testing.T) {
	initClient(nil)
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cardCmdFlags(cmd)

	name := "test_secret"
	err := cmd.Flags().Set("name", name)
	require.NoError(t, err)

	cardNumber := "1234567890123456"
	err = cmd.Flags().Set("number", cardNumber)
	require.NoError(t, err)

	date := "12/24"
	err = cmd.Flags().Set("date", date)
	require.NoError(t, err)

	code := "123"
	err = cmd.Flags().Set("code", code)
	require.NoError(t, err)

	holder := "John Doe"
	err = cmd.Flags().Set("holder", holder)
	require.NoError(t, err)

	token_cache.InitTokenStorage("./test/token.txt")

	card := vaulttypes.Card{
		Number:       cardNumber,
		ExpiryDate:   date,
		SecurityCode: code,
		Holder:       holder,
	}

	// Encrypt the card details.
	mockCipher := mocks.NewCipher(t)
	mockCipher.On("Encrypt", card).Return(nil, errors.New("Internal server error"))
	initCipher(mockCipher)

	err = createCardCmdRunE(cmd, []string{})
	assert.Error(t, err)
}
