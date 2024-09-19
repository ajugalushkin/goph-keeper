package text

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

var keepUpdateText = &cobra.Command{
	Use:   "text",
	Short: "Update text secret",
	RunE:  keeperUpdateTextCmdRunE,
}

var client app.KeeperClient

func NewCommand() *cobra.Command {
	return keepUpdateText
}

// keeperUpdateTextCmdRunE is the entry point for the "text" command in the "keep update" command group.
// It updates a text secret in the Goph-Keeper vault.
//
// Parameters:
// - cmd: A pointer to the cobra.Command object representing the "text" command.
// - args: An array of strings containing any additional arguments passed to the command.
//
// Returns:
//   - An error if any error occurs during the execution of the command.
//     If no error occurs, it returns nil.
func keeperUpdateTextCmdRunE(cmd *cobra.Command, args []string) error {
	const op = "keep update text"
	log := logger.GetLogger().With("op", op)

	// Read the secret name from command-line flags
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ",
			slog.String("error", err.Error()))
		return fmt.Errorf("error reading secret name")
	}
	if name == "" {
		return fmt.Errorf("secret name cannot be empty")
	}

	// Read the text data from command-line flags
	data, err := cmd.Flags().GetString("data")
	if err != nil {
		log.Error("Error reading text data: ",
			slog.String("error", err.Error()))
		return fmt.Errorf("error reading text data")
	}
	if data == "" {
		return fmt.Errorf("text data cannot be empty")
	}

	// Create a Text secret object with the provided data
	text := vaulttypes.Text{
		Data: data,
	}

	// Encrypt the secret data
	content, err := secret.NewCryptographer().Encrypt(text)
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
	// Send an update request to the vault
	resp, err := client.UpdateItem(context.Background(), &v1.UpdateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to update secret: ",
			slog.String("error", err.Error()))
		return err
	}

	// Print a success message
	fmt.Printf("Secret %s version %v updated successfully\n", resp.GetName(), resp.GetVersion())
	return nil
}

// updateTextCmdFlags sets up command-line flags for updating text secrets.
//
// The function takes a pointer to a cobra.Command object as a parameter.
// It adds two flags to the command:
// - "name": A string flag representing the name of the secret to be updated.
// - "data": A string flag representing the text data to be stored in the secret.
//
// The flags are configured with default values of "" and descriptive help messages.
func updateTextCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Secret name")
	cmd.Flags().String("data", "", "Text data")
}

func init() {
	updateTextCmdFlags(keepUpdateText)
}

func initClient(newClient app.KeeperClient) {
	client = newClient
}
