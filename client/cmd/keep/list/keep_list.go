package list

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/ajugalushkin/goph-keeper/client/secret"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

var client app.KeeperClient

// NewCommand creates a new Cobra command for listing secrets.
// The command is configured to use the "list" subcommand, with a short description "List secrets".
// When the command is executed, the keepListRun function is called.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List secrets",
		RunE:  keepListRunE,
	}
}

// keepListRun is a command handler function that lists secrets stored in the goph-keeper service.
// It retrieves the access token_cache, establishes a connection to the goph-keeper service, sends a request to list secrets,
// decrypts the received secrets, and prints them to the console.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
// - args: A slice of strings representing command-line arguments.
//
// Return:
// This function does not return any value.
func keepListRunE(cmd *cobra.Command, args []string) error {
	const op = "keep_get"
	log := logger.GetLogger().With("op", op)

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

	resp, err := client.ListItems(context.Background(), &v1.ListItemsRequestV1{})
	if err != nil {
		log.Error("Failed to list secret: ", slog.String("error", err.Error()))
		return err
	}

	for _, info := range resp.GetSecrets() {
		newSecret, err := secret.Decrypt(info.GetContent())
		if err != nil {
			log.Error("Failed to decrypt secret: ", slog.String("error", err.Error()))
			return err
		}

		fmt.Printf("%s\n", newSecret)
	}
	return nil
}

func initClient(newClient app.KeeperClient) {
	client = newClient
}
