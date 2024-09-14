package get

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/get/bin"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/secret"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

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
	const op = "keep_get"

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get secret",
		Run:   keepGetRun,
	}

	cmd.AddCommand(bin.NewCommand())

	cmd.Flags().String("name", "", "Secret name")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	return cmd
}

// keepGetRun is the function that handles the "get" command for the keep command.
// It retrieves a secret from the goph-keeper service based on the provided secret name.
//
// Parameters:
// - cmd: The Cobra command object representing the "get" command.
// - args: Additional arguments provided to the command.
//
// Return:
// - None.
func keepGetRun(cmd *cobra.Command, args []string) {
	const op = "keep_get"
	log := logger.GetInstance().Log.With("op", op)

	// Read the secret name from the command flags.
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ", slog.String("error", err.Error()))
	}

	// Load the authentication token_cache from the token_cache storage.
	token, err := token_cache.GetInstance().Load()
	if err != nil {
		return
	}

	// Get the client configuration for connecting to the goph-keeper service.
	cfg := config.GetInstance().Config.Client

	// Create a new keeper client using the provided configuration and authentication token_cache.
	keeperClient := keeper.NewKeeperClient(keeper.GetKeeperConnection(log, cfg.Address, token))

	// Request the secret from the goph-keeper service using the secret name.
	resp, err := keeperClient.GetItem(context.Background(), &v1.GetItemRequestV1{
		Name: name,
	})
	if err != nil {
		log.Error("Failed to create secret: ", slog.String("error", err.Error()))
	}

	var respSecret vaulttypes.Vault
	if resp == nil {
		log.Error("Secret not found")
		return
	} else {
		// Decrypt the retrieved secret content.
		respSecret, err = secret.DecryptSecret(resp.GetContent())
		if err != nil {
			log.Error("Failed to decrypt secret: ", slog.String("error", err.Error()))
		}
	}

	// Print the decrypted secret to the console.
	fmt.Printf("%s\n", respSecret)
}
