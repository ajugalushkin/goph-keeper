package get

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/get/bin"
	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/ajugalushkin/goph-keeper/client/secret"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

var keepGet = &cobra.Command{
	Use:   "get",
	Short: "Get secret",
	RunE:  keepGetRunE,
}

var client app.KeeperClient
var cipher secret.Cipher

// NewCommand creates a new Cobra command for retrieving a secret from the goph-keeper service.
// The command is named "get" and is used to interact with the secret storage system.
// It accepts a required flag "--name" to specify the name of the secret to retrieve.
//
// The command also includes a subcommand for handling binary data, which is added using the bin.NewCommand() function.
//
// If the "--name" flag is not provided, an error message is logged using the slog package.
//
// The function returns a pointer to the created Cobra command.
func NewCommand() *cobra.Command {
	return keepGet
}

// keepGetRunE is the entry point for the "get" command in the goph-keeper client.
// It retrieves a secret from the goph-keeper service based on the provided command-line arguments.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
// - args: An array of strings representing the command-line arguments.
//
// Returns:
// - An error if any error occurs during the execution of the command.
func keepGetRunE(cmd *cobra.Command, args []string) error {
	const op = "keep_get"
	log := logger.GetLogger().With("op", op)

	// Read the secret name from the command flags.
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ", slog.String("error", err.Error()))
		return err
	}
	if name == "" {
		log.Error("Secret name is required")
		return fmt.Errorf("secret name is required")
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
	// Request the secret from the goph-keeper service using the secret name.
	resp, err := client.GetItem(context.Background(), &v1.GetItemRequestV1{
		Name: name,
	})
	if err != nil {
		log.Error("Failed to create secret: ", slog.String("error", err.Error()))
		return err
	}

	// Decrypt the retrieved secret content.
	if cipher == nil {
		cipher = secret.NewCryptographer()
	}
	// Decrypt the retrieved secret content using the provided cipher.
	respSecret, err := cipher.Decrypt(resp.GetContent())
	if err != nil {
		return err
	}

	// Print the decrypted secret to the console.
	fmt.Printf("%s\n", respSecret.String())
	return nil
}
func getCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Secret name")
}

// init initializes the command-line flags for the "get" command and adds a subcommand for handling binary data.
//
// Parameters:
// - keepGet: A pointer to the Cobra command object representing the "get" command.
//
// Returns:
// - None
func init() {
	getCmdFlags(keepGet)
	keepGet.AddCommand(bin.NewCommand())
}

func initCipher(newCipher secret.Cipher) {
	cipher = newCipher
}
