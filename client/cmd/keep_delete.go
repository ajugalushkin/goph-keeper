package cmd

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"

	"log/slog"

	"github.com/spf13/cobra"
)

// keepDeleteCmd represents the delete command
var keepDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Run:   keepDeleteCmdRun,
}

// init initializes the delete command for the keep command.
// It adds a flag for the secret name and marks it as required.
// If an error occurs while setting the flag, it logs the error using the provided logger.
func init() {
    const op = "keep_create_card"
    keepCmd.AddCommand(keepDeleteCmd)

    // Add a flag for the secret name.
    keepDeleteCmd.Flags().String("name", "", "Secret name")

    // Mark the secret name flag as required.
    if err := keepDeleteCmd.MarkFlagRequired("name"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }
}

// keepDeleteCmdRun handles the deletion of a secret from the goph-keeper service.
// It reads the secret name from the command-line flags, retrieves the user's token,
// establishes a connection to the goph-keeper service, and sends a delete request.
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

    // Retrieve the user's token.
    token, err := tokenStorage.Load()
    if err != nil {
        return
    }

    // Retrieve the goph-keeper client configuration.
    cfg := config.GetInstance().Config.Client

    // Establish a connection to the goph-keeper service.
    keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

    // Send a delete request to the goph-keeper service.
    resp, err := keeperClient.DeleteItem(context.Background(), &keeperv1.DeleteItemRequestV1{
        Name: name,
    })
    if err != nil {
        log.Error("Failed to delete secret", slog.String("error", err.Error()))
        return
    }

    // Print a success message.
    fmt.Printf("Secret %s deleted successfully\n", resp.GetName())
}
