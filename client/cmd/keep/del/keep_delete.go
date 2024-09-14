package del

import (
	"context"
	"fmt"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"

	"log/slog"

	"github.com/spf13/cobra"
)

var keepDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a secret from the goph-keeper service",
	RunE:  keepDeleteCmdRunE,
}

var client app.KeeperClient

// NewCommand creates a new Cobra command for deleting a secret from the goph-keeper service.
// The command is configured to accept a "name" flag, which specifies the secret to be deleted.
// The command also ensures that the "name" flag is required.
// Upon execution, the command calls the keepDeleteCmdRun function to handle the deletion process.
func NewCommand() *cobra.Command {
	return keepDelete
}

// keepDeleteCmdRunE is the main function for handling the deletion of a secret from the goph-keeper service.
// It reads the secret name from the command-line flags, initializes the Keeper client if not already done,
// sends a delete request to the goph-keeper service, and prints a success message upon successful deletion.
//
// Parameters:
// - cmd: The Cobra command object.
// - args: Additional command-line arguments.
//
// Returns:
// - An error if any error occurs during the deletion process.
func keepDeleteCmdRunE(cmd *cobra.Command, args []string) error {
	const op = "keep_delete"
	log := logger.GetLogger().With("op", op)

	// Read the secret name from the command-line flags.
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name")
		return fmt.Errorf("error reading secret name")
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
	// Send a del request to the goph-keeper service.
	resp, err := client.DeleteItem(context.Background(), &keeperv1.DeleteItemRequestV1{
		Name: name,
	})
	if err != nil {
		log.Error("Failed to del secret", slog.String("error", err.Error()))
		return err
	}

	// Print a success message.
	fmt.Printf("Secret %s deleted successfully\n", resp.GetName())
	return nil
}

func deleteCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Name of the secret to be deleted")
}

func init() {
	deleteCmdFlags(keepDelete)
}

func initClient(newClient app.KeeperClient) {
	client = newClient
}
