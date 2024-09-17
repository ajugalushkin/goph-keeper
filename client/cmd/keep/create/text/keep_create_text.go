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

var keepCreateText = &cobra.Command{
	Use:   "text",
	Short: "Create text secret",
	RunE:  keepCreateTextCmdRunE,
}

var client app.KeeperClient

// NewCommand creates a cobra.Command for creating a text secret.
// The command is responsible for handling the execution of the "text" subcommand of the "keep create" command.
// It reads the required flags "name" and "data", creates a text secret, encrypts it, and sends a request to the goph-keeper server to store the secret.
// If any error occurs during the process, it logs the error using the slog package.
//
// Parameters:
// - None
//
// Return:
// - A pointer to the cobra.Command object representing the "text" subcommand.
func NewCommand() *cobra.Command {
	return keepCreateText
}

// keepCreateTextCmdRunE is the entry point for the "text" subcommand of the "keep create" command.
// It handles the execution of creating a text secret by reading the required flags, encrypting the secret,
// and sending a request to the goph-keeper server to store the secret.
//
// Parameters:
// - cmd: A pointer to the cobra.Command object representing the "text" subcommand.
// - args: An array of strings containing any additional arguments passed to the command.
//
// Return:
// - An error if any error occurs during the process, or nil if the operation is successful.
func keepCreateTextCmdRunE(cmd *cobra.Command, args []string) error {
	const op = "keep_create_text"
	log := logger.GetLogger().With("op", op)

	// Read the secret name from the command flags
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ",
			slog.String("error", err.Error()))
		return err
	}

	// Read the text data from the command flags
	data, err := cmd.Flags().GetString("data")
	if err != nil {
		log.Error("Error reading text data: ",
			slog.String("error", err.Error()))
		return err
	}

	// Create a Text secret object
	text := vaulttypes.Text{
		Data: data,
	}

	// Encrypt the secret
	content, err := secret.EncryptSecret(text)
	if err != nil {
		log.Error("Failed to encrypt secret: ",
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
	// Send a request to the goph-keeper server to create the secret
	resp, err := client.CreateItem(context.Background(), &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to create secret: ", slog.String("error", err.Error()))
		return err
	}

	// Check if the response is nil
	if resp == nil {
		log.Error("Nil response received from Keeper server")
		return fmt.Errorf("nil response received from Keeper server")
	}

	// Print a success message
	fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())

	return nil
}

// textCmdFlags sets up the command line flags for the "text" subcommand of the "keep create" command.
// It adds two flags: "name" and "data". The "name" flag is used to specify the name of the secret,
// while the "data" flag is used to provide the text data for the secret.
//
// Parameters:
// - cmd: A pointer to the cobra.Command object representing the "text" subcommand.
//
// Return:
// - None
func textCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Secret name")
	cmd.Flags().String("data", "", "Text data")
}

func init() {
	textCmdFlags(keepCreateText)
}
