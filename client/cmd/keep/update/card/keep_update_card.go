package card

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/secret"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// NewCommand returns a Cobra command for updating a card secret in the Keeper service.
// The command accepts flags for the secret name, card number, expiry date, security code, and holder.
// It then creates a Card struct, encrypts the secret, and sends an update request to the Keeper service.
// If any error occurs during the process, it logs the error and returns.
func NewCommand() *cobra.Command {
	const op = "keep update card"

	cmd := &cobra.Command{
		Use:   "card",
		Short: "Update card secret",
		Run:   keeperUpdateCardCmdRun,
	}

	// Define and mark required flags for the command
	cmd.Flags().String("name", "", "Secret name")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	cmd.Flags().String("number", "", "Card number")
	if err := cmd.MarkFlagRequired("number"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	cmd.Flags().String("date", "", "Card expiry date")
	if err := cmd.MarkFlagRequired("date"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	cmd.Flags().String("code", "", "Card security code")
	if err := cmd.MarkFlagRequired("code"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	cmd.Flags().String("holder", "", "Card holder")
	if err := cmd.MarkFlagRequired("holder"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	return cmd
}

// keeperUpdateCardCmdRun is the function that handles the "card" command for updating a card secret.
// It reads the required flags, creates a Card struct, encrypts the secret, and sends an update request to the Keeper service.
// If any error occurs during the process, it logs the error and returns.
//
// Parameters:
// - cmd: The Cobra command object.
// - args: Additional command-line arguments.
func keeperUpdateCardCmdRun(cmd *cobra.Command, args []string) {
	const op = "keep update card"
	log := logger.GetLogger().With("op", op)

	// Read the required flags
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ",
			slog.String("error", err.Error()))
		return
	}

	number, err := cmd.Flags().GetString("number")
	if err != nil {
		log.Error("Error reading card number: ",
			slog.String("error", err.Error()))
		return
	}

	date, err := cmd.Flags().GetString("date")
	if err != nil {
		log.Error("Error reading card expiry date: ",
			slog.String("error", err.Error()))
		return
	}

	code, err := cmd.Flags().GetString("code")
	if err != nil {
		log.Error("Error reading card security code: ",
			slog.String("error", err.Error()))
		return
	}

	holder, err := cmd.Flags().GetString("holder")
	if err != nil {
		log.Error("Error reading card holder: ",
			slog.String("error", err.Error()))
		return
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
		return
	}

	// Load the token_cache
	token, err := token_cache.GetToken().Load()
	if err != nil {
		return
	}

	// Get the Keeper client configuration
	cfg := config.GetConfig().Client
	keeperClient := keeper.NewKeeperClient(keeper.GetKeeperConnection(log, cfg.Address, token))

	// Send the update request to the Keeper service
	resp, err := keeperClient.UpdateItem(context.Background(), &v1.UpdateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to update secret: ",
			slog.String("error", err.Error()))
		return
	}

	// Print the success message
	fmt.Printf("Secret %s version %v updated successfully\n", resp.GetName(), resp.GetVersion())
}
