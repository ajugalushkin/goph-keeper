package cmd

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// keepUpdateCardCmd represents the card command
var keepUpdateCardCmd = &cobra.Command{
	Use:   "card",
	Short: "Update card secret",
	Run:   keeperUpdateCardCmdRun,
}

// init initializes the "card" command for updating card secret.
// It adds the command to the "keep update" command and sets up required flags.
// The flags include:
// - "name": The name of the secret to be updated.
// - "number": The card number.
// - "date": The card expiry date.
// - "code": The card security code.
// - "holder": The card holder.
// If any of the required flags are missing, an error is logged.
func init() {
    const op = "keep update card"
    keepUpdateCmd.AddCommand(keepUpdateCardCmd)

    keepUpdateCardCmd.Flags().String("name", "", "Secret name")
    if err := keepUpdateCardCmd.MarkFlagRequired("name"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }
    keepUpdateCardCmd.Flags().String("number", "", "Card number")
    if err := keepUpdateCardCmd.MarkFlagRequired("number"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }
    keepUpdateCardCmd.Flags().String("date", "", "Card expiry date")
    if err := keepUpdateCardCmd.MarkFlagRequired("date"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }
    keepUpdateCardCmd.Flags().String("code", "", "Card security code")
    if err := keepUpdateCardCmd.MarkFlagRequired("code"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }
    keepUpdateCardCmd.Flags().String("holder", "", "Card holder")
    if err := keepUpdateCardCmd.MarkFlagRequired("holder"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }
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
    log := logger.GetInstance().Log.With("op", op)

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
    content, err := encryptSecret(card)
    if err != nil {
        log.Error("Failed to encrypt secret: ",
            slog.String("error", err.Error()))
        return
    }

    // Load the token
    token, err := tokenStorage.Load()
    if err != nil {
        return
    }

    // Get the Keeper client configuration
    cfg := config.GetInstance().Config.Client
    keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

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
