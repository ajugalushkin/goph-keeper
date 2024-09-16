package card

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/ajugalushkin/goph-keeper/client/secret"
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

var keepUpdateCard = &cobra.Command{
	Use:   "card",
	Short: "Update card secret",
	RunE:  keeperUpdateCardCmdRunE,
}

var client app.KeeperClient

// NewCommand returns a Cobra command for updating a card secret in the Keeper service.
// The command accepts flags for the secret name, card number, expiry date, security code, and holder.
// It then creates a Card struct, encrypts the secret, and sends an update request to the Keeper service.
// If any error occurs during the process, it logs the error and returns.
func NewCommand() *cobra.Command {
	return keepUpdateCard
}

// keeperUpdateCardCmdRun is the function that handles the "card" command for updating a card secret.
// It reads the required flags, creates a Card struct, encrypts the secret, and sends an update request to the Keeper service.
// If any error occurs during the process, it logs the error and returns.
//
// Parameters:
// - cmd: The Cobra command object.
// - args: Additional command-line arguments.
func keeperUpdateCardCmdRunE(cmd *cobra.Command, args []string) error {
	const op = "keep update card"
	log := logger.GetLogger().With("op", op)

	// Read the required flags
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ",
			slog.String("error", err.Error()))
		return fmt.Errorf("required flag(s) \"name\" not set")
	}
	if name == "" {
		return fmt.Errorf("invalid secret name")
	}

	number, err := cmd.Flags().GetString("number")
	if err != nil {
		log.Error("Error reading card number: ",
			slog.String("error", err.Error()))
		return fmt.Errorf("error reading card number")
	}
	if number == "" {
		return fmt.Errorf("invalid card number")
	}

	date, err := cmd.Flags().GetString("date")
	if err != nil {
		log.Error("Error reading card expiry date: ",
			slog.String("error", err.Error()))
		return err
	}
	if date == "" {
		return fmt.Errorf("invalid card expiry date")
	}

	code, err := cmd.Flags().GetString("code")
	if err != nil {
		log.Error("Error reading card security code: ",
			slog.String("error", err.Error()))
		return err
	}
	if code == "" {
		return fmt.Errorf("invalid card security code")
	}

	holder, err := cmd.Flags().GetString("holder")
	if err != nil {
		log.Error("Error reading card holder: ",
			slog.String("error", err.Error()))
		return err
	}
	if holder == "" {
		return fmt.Errorf("invalid card holder")
	}

	// Create a Card struct
	card := vaulttypes.Card{
		Number:       number,
		ExpiryDate:   date,
		SecurityCode: code,
		Holder:       holder,
	}

	// Encrypt the secret
	content, err := secret.EncryptSecret(card)
	if err != nil {
		log.Error("Failed to secret secret: ",
			slog.String("error", err.Error()))
		return err
	}

	// If the Keeper client is not initialized, load the authentication token_cache from storage and create a new client.
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
	// Send the update request to the Keeper service
	resp, err := client.UpdateItem(context.Background(), &v1.UpdateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to update secret: ",
			slog.String("error", err.Error()))
		return err
	}

	// Print the success message
	fmt.Printf("Secret %s version %v updated successfully\n", resp.GetName(), resp.GetVersion())
	return nil
}

// updateCardCmdFlags sets up the command-line flags for updating a card secret in the Keeper service.
// It accepts five flags: name, number, date, code, and holder.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
//
// The function adds the following flags to the command:
// - "name": The name of the secret to be updated.
// - "number": The card number.
// - "date": The card expiry date.
// - "code": The card security code.
// - "holder": The card holder's name.
func updateCardCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Secret name")
	cmd.Flags().String("number", "", "Card number")
	cmd.Flags().String("date", "", "Card expiry date")
	cmd.Flags().String("code", "", "Card security code")
	cmd.Flags().String("holder", "", "Card holder")
}

func init() {
	updateCardCmdFlags(keepUpdateCard)
}

func initClient(newClient app.KeeperClient) {
	client = newClient
}
