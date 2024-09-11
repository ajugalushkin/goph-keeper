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

// textCmd represents the text command
var keepUpdateTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Update text secret",
	Run:   keeperUpdateTextCmdRun,
}

// init initializes the "keep update text" command.
// This command is responsible for updating a text secret in the Goph-Keeper vault.
// It reads the secret name and data from command-line flags, encrypts the data, and sends an update request to the vault.
// If any errors occur during the process, they are logged and the function returns.
// If the secret is successfully updated, a success message is printed.
func init() {
    const op = "keep update text"
    keepUpdateCmd.AddCommand(keepUpdateTextCmd)

    // Add a flag for the secret name. The flag is required.
    keepUpdateTextCmd.Flags().String("name", "", "Secret name")
    if err := keepUpdateTextCmd.MarkFlagRequired("name"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }

    // Add a flag for the text data. The flag is required.
    keepUpdateTextCmd.Flags().String("data", "", "Text data")
    if err := keepUpdateTextCmd.MarkFlagRequired("data"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }
}

// keeperUpdateTextCmdRun is responsible for updating a text secret in the Goph-Keeper vault.
// It reads the secret name and data from command-line flags, encrypts the data, and sends an update request to the vault.
// If any errors occur during the process, they are logged and the function returns.
// If the secret is successfully updated, a success message is printed.
func keeperUpdateTextCmdRun(cmd *cobra.Command, args []string) {
	const op = "keep update text"
	log := logger.GetInstance().Log.With("op", op)

	// Read the secret name from command-line flags
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ",
			slog.String("error", err.Error()))
		return
	}

	// Read the text data from command-line flags
	data, err := cmd.Flags().GetString("data")
	if err != nil {
		log.Error("Error reading text data: ",
			slog.String("error", err.Error()))
		return
	}

	// Create a Text secret object with the provided data
	text := vaulttypes.Text{
		Data: data,
	}

	// Encrypt the secret data
	content, err := encryptSecret(text)
	if err != nil {
		log.Error("Failed to encrypt secret: ",
			slog.String("error", err.Error()))
		return
	}

	// Load the authentication token from storage
	token, err := tokenStorage.Load()
	if err != nil {
		return
	}

	// Get a connection to the Goph-Keeper vault
	cfg := config.GetInstance().Config.Client
	keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

	// Send an update request to the vault
	resp, err := keeperClient.UpdateItem(context.Background(), &v1.UpdateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to update secret: ",
			slog.String("error", err.Error()))
		return
	}

	// Print a success message
	fmt.Printf("Secret %s version %v updated successfully\n", resp.GetName(), resp.GetVersion())
}
