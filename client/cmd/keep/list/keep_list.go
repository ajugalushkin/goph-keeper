package list

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
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// NewCommand creates a new Cobra command for listing secrets.
// The command is configured to use the "list" subcommand, with a short description "List secrets".
// When the command is executed, the keepListRun function is called.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List secrets",
		Run:   keepListRun,
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
func keepListRun(cmd *cobra.Command, args []string) {
	const op = "keep_get"
	log := logger.GetInstance().Log.With("op", op)

	token, err := token_cache.GetInstance().Load()
	if err != nil {
		return
	}

	cfg := config.GetInstance().Config.Client
	keeperClient := keeper.NewKeeperClient(keeper.GetKeeperConnection(log, cfg.Address, token))
	resp, err := keeperClient.ListItems(context.Background(), &v1.ListItemsRequestV1{})
	if err != nil {
		log.Error("Failed to list secret: ", slog.String("error", err.Error()))
	}

	for _, info := range resp.GetSecrets() {
		newSecret, err := secret.DecryptSecret(info.GetContent())
		if err != nil {
			log.Error("Failed to decrypt secret: ", slog.String("error", err.Error()))
		}

		fmt.Printf("%s\n", newSecret)
	}
}
