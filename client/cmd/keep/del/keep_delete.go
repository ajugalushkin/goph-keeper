package del

import (
	"context"
	"fmt"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"

	"log/slog"

	"github.com/spf13/cobra"
)

// NewCommand creates a new Cobra command for deleting a secret from the goph-keeper service.
// The command is configured to accept a "name" flag, which specifies the secret to be deleted.
// The command also ensures that the "name" flag is required.
// Upon execution, the command calls the keepDeleteCmdRun function to handle the deletion process.
func NewCommand() *cobra.Command {
	const op = "keep_delete"

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a secret from the goph-keeper service",
		Run:   keepDeleteCmdRun,
	}

	// Add a flag for the secret name.
	cmd.Flags().String("name", "", "Name of the secret to be deleted")

	// Mark the secret name flag as required.
	if err := cmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	return cmd
}

// keepDeleteCmdRun handles the deletion of a secret from the goph-keeper service.
// It reads the secret name from the command-line flags, retrieves the user's token_cache,
// establishes a connection to the goph-keeper service, and sends a del request.
// If the deletion is successful, it prints a success message.
//
// Parameters:
// - cmd: The Cobra command object.
// - args: Additional command-line arguments.
//
// Returns:
// - None.
func keepDeleteCmdRun(cmd *cobra.Command, args []string) {
	const op = "keep_delete"
	log := logger.GetInstance().Log.With("op", op)

	// Read the secret name from the command-line flags.
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name", slog.String("error", err.Error()))
	}

	// Retrieve the user's token_cache.
	token, err := token_cache.GetInstance().Load()
	if err != nil {
		return
	}

	// Retrieve the goph-keeper client configuration.
	cfg := config.GetInstance().Config.Client

	// Establish a connection to the goph-keeper service.
	keeperClient := keeper.NewKeeperClient(keeper.GetKeeperConnection(log, cfg.Address, token))

	// Send a del request to the goph-keeper service.
	resp, err := keeperClient.DeleteItem(context.Background(), &keeperv1.DeleteItemRequestV1{
		Name: name,
	})
	if err != nil {
		log.Error("Failed to del secret", slog.String("error", err.Error()))
		return
	}

	// Print a success message.
	fmt.Printf("Secret %s deleted successfully\n", resp.GetName())
}
