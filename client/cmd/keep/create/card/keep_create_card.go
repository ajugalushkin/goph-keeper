package card

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/secret"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"

	"github.com/spf13/cobra"
)

// NewCommand creates a Cobra command for creating a card secret.
// The command accepts flags for specifying the secret name, card number, expiry date,
// security code, and card holder. It then encrypts the card details, sends a request to
// the Keeper server to create the secret, and prints a success message with the created
// secret's name and version.
func NewCommand() *cobra.Command {
	const op = "keep_create_card"

	cmd := &cobra.Command{
		Use:   "card",
		Short: "Create card secret",
		Run:   createCardCmdRun,
	}

	// "name" flag is used to specify the secret name.
	cmd.Flags().String("name", "", "Secret name")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	// "number" flag is used to specify the card number.
	cmd.Flags().String("number", "", "Card number")
	if err := cmd.MarkFlagRequired("number"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	// "date" flag is used to specify the card expiry date.
	cmd.Flags().String("date", "", "Card expiry date")
	if err := cmd.MarkFlagRequired("date"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	// "code" flag is used to specify the card security code.
	cmd.Flags().String("code", "", "Card security code")
	if err := cmd.MarkFlagRequired("code"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	// "holder" flag is used to specify the card holder.
	cmd.Flags().String("holder", "", "Card holder")
	if err := cmd.MarkFlagRequired("holder"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	return cmd
}

// createCardCmdRun is responsible for handling the execution of the "card" command.
// It reads the required flags for creating a card secret, encrypts the card details,
// and sends a request to the Keeper server to create the secret.
//
// Parameters:
// - cmd: The Cobra command object.
// - args: Additional command-line arguments.
func createCardCmdRun(cmd *cobra.Command, args []string) {
	const op = "keep_create_card"
	log := logger.GetLogger().With("op", op)

	// Read the required flags for creating a card secret.
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

	// Create a Card object from the provided details.
	card := vaulttypes.Card{
		Number:       number,
		ExpiryDate:   date,
		SecurityCode: code,
		Holder:       holder,
	}

	// Encrypt the card details.
	content, err := secret.EncryptSecret(card)
	if err != nil {
		log.Error("Failed to secret secret: ",
			slog.String("error", err.Error()))
		return
	}

	// Load the authentication token_cache from storage.
	token, err := token_cache.GetToken().Load()
	if err != nil {
		return
	}

	// Create a new Keeper client using the provided configuration and token_cache.
	cfg := config.GetConfig().Client
	keeperClient := keeper.NewKeeperClient(keeper.GetKeeperConnection(log, cfg.Address, token))

	// Send a request to the Keeper server to create the secret.
	resp, err := keeperClient.CreateItem(context.Background(), &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to create secret: ", slog.String("error", err.Error()))
		return
	}
	if resp == nil {
		log.Error("No response received from Keeper server")
		return
	}

	// Print a success message with the created secret's name and version.
	fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())
}
