package cmd

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"log/slog"

	"github.com/spf13/cobra"
)

// keepCreateCardCmd represents the card command
var keepCreateCardCmd = &cobra.Command{
	Use:   "card",
	Short: "Create card secret",
	Run:   createCardCmdRun,
}

// init initializes the "card" command for creating a card secret.
// It sets up required flags for the card details and logs any errors during flag setup.
func init() {
	const op = "keep_create_card"
	keepCreateCmd.AddCommand(keepCreateCardCmd)

	// "name" flag is used to specify the secret name.
	keepCreateCardCmd.Flags().String("name", "", "Secret name")
	if err := keepCreateCardCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	// "number" flag is used to specify the card number.
	keepCreateCardCmd.Flags().String("number", "", "Card number")
	if err := keepCreateCardCmd.MarkFlagRequired("number"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	// "date" flag is used to specify the card expiry date.
	keepCreateCardCmd.Flags().String("date", "", "Card expiry date")
	if err := keepCreateCardCmd.MarkFlagRequired("date"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	// "code" flag is used to specify the card security code.
	keepCreateCardCmd.Flags().String("code", "", "Card security code")
	if err := keepCreateCardCmd.MarkFlagRequired("code"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	// "holder" flag is used to specify the card holder.
	keepCreateCardCmd.Flags().String("holder", "", "Card holder")
	if err := keepCreateCardCmd.MarkFlagRequired("holder"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
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
	log := logger.GetInstance().Log.With("op", op)

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
	content, err := encryptSecret(card)
	if err != nil {
		log.Error("Failed to encrypt secret: ",
			slog.String("error", err.Error()))
		return
	}

	// Load the authentication token from storage.
	token, err := tokenStorage.Load()
	if err != nil {
		return
	}

	// Create a new Keeper client using the provided configuration and token.
	cfg := config.GetInstance().Config.Client
	keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

	// Send a request to the Keeper server to create the secret.
	resp, err := keeperClient.CreateItem(context.Background(), &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to create secret: ", slog.String("error", err.Error()))
	}

	// Print a success message with the created secret's name and version.
	fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())
}
