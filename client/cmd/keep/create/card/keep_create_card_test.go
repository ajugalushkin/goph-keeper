package card

import (
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// Fails to read the secret name flag
func TestCreateCardCmdRunE_FailsToReadSecretNameFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("number", "1234567890123456", "card number")
	cmd.Flags().String("date", "12/24", "expiry date")
	cmd.Flags().String("code", "123", "security code")
	cmd.Flags().String("holder", "John Doe", "card holder")

	err := createCardCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

// Fails to read the card number flag
func TestCreateCardCmdRunE_FailsToReadCardNumberFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name of the secret")
	cmd.Flags().String("date", "12/24", "expiry date")
	cmd.Flags().String("code", "123", "security code")
	cmd.Flags().String("holder", "John Doe", "card holder")

	err := createCardCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatal("Expected an error but got nil")
	}
}

// Fails to read the card expiry date flag
func TestCreateCardCmdRunE_FailsToReadExpiryDateFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name of the secret")
	cmd.Flags().String("number", "1234567890123456", "card number")
	cmd.Flags().String("date", "", "expiry date") // Simulating failure to read expiry date flag
	cmd.Flags().String("code", "123", "security code")
	cmd.Flags().String("holder", "John Doe", "card holder")

	err := createCardCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatal("Expected an error but got nil")
	}
}
