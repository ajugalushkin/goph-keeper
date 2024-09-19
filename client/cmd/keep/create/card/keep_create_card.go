package card

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/ajugalushkin/goph-keeper/client/secret"
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"

	"github.com/spf13/cobra"
)

var keepCreateCard = NewCommand()

var client app.KeeperClient

// NewCommand creates a Cobra command for creating a card secret.
// The command accepts flags for specifying the secret name, card number, expiry date,
// security code, and card holder. It then encrypts the card details, sends a request to
// the Keeper server to create the secret, and prints a success message with the created
// secret's name and version.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "card",
		Short: "Create card secret",
		RunE:  createCardCmdRunE,
	}
}

// createCardCmdRunE is the entry point for creating a card secret in the Keeper vault.
// It reads the required flags for creating a card secret, encrypts the card details,
// sends a request to the Keeper server to create the secret, and prints a success message.
//
// Parameters:
// - cmd: The Cobra command object.
// - args: Additional command-line arguments.
//
// Returns:
// - An error if any error occurs during the process, or nil if the operation is successful.
func createCardCmdRunE(cmd *cobra.Command, args []string) error {
	const op = "keep_create_card"
	log := logger.GetLogger().With("op", op)

	// Read the required flags for creating a card secret.
	name, err := cmd.Flags().GetString("name")
	if err != nil || name == "" {
		if err != nil {
			log.Error("Error reading secret name: ", slog.String("error", err.Error()))
		}
		return fmt.Errorf("name is required")
	}

	number, err := cmd.Flags().GetString("number")
	if err != nil || number == "" {
		if err != nil {
			log.Error("Error reading card number: ", slog.String("error", err.Error()))
		}
		return fmt.Errorf("card number is required")
	}

	date, err := cmd.Flags().GetString("date")
	if err != nil || date == "" {
		if err != nil {
			log.Error("Error reading card expiry date: ", slog.String("error", err.Error()))
		}
		return fmt.Errorf("expiry date is required")
	}

	code, err := cmd.Flags().GetString("code")
	if err != nil || code == "" {
		if err != nil {
			log.Error("Error reading card security code: ",
				slog.String("error", err.Error()))
		}
		return fmt.Errorf("security code is required")
	}

	holder, err := cmd.Flags().GetString("holder")
	if err != nil || holder == "" {
		if err != nil {
			log.Error("Error reading card holder: ",
				slog.String("error", err.Error()))
		}
		return fmt.Errorf("card holder is required")
	}

	// Create a Card object from the provided details.
	card := vaulttypes.Card{
		Number:       number,
		ExpiryDate:   date,
		SecurityCode: code,
		Holder:       holder,
	}

	// Encrypt the card details.
	content, err := secret.NewCryptographer().Encrypt(card)
	if err != nil {
		log.Error("Failed to secret secret: ",
			slog.String("error", err.Error()))
		return err
	}

	if client == nil {
		// Load the authentication token_cache from storage.
		token, err := token_cache.GetToken().Load()
		if err != nil {
			return err
		}
		// Create a new Keeper client using the provided configuration and token_cache.
		cfg := config.GetConfig().Client
		client = keeper.NewKeeperClient(keeper.GetKeeperConnection(log, cfg.Address, token))
	}

	// Send a request to the Keeper server to create the secret.
	resp, err := client.CreateItem(context.Background(), &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to create secret: ", slog.String("error", err.Error()))
		return err
	}
	if resp == nil {
		log.Error("No response received from Keeper server")
		return err
	}

	// Print a success message with the created secret's name and version.
	fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())
	return nil
}

// cardCmdFlags sets up the command-line flags for creating a card secret in the Keeper vault.
// It accepts five flags: name, number, date, code, and holder. These flags are used to specify
// the secret name, card number, expiry date, security code, and card holder, respectively.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
//
// The function does not return any value. It sets up the command-line flags for the given command.
func cardCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Secret name")
	cmd.Flags().String("number", "", "Card number")
	cmd.Flags().String("date", "", "Card expiry date")
	cmd.Flags().String("code", "", "Card security code")
	cmd.Flags().String("holder", "", "Card holder")
}

func init() {
	cardCmdFlags(keepCreateCard)
}

func initClient(newClient app.KeeperClient) {
	client = newClient
}
