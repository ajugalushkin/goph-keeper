package cmd

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// getCmd represents the get command
var keepGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secret",
	Run:   keepGetRun,
}

// init initializes the get command for the keep command.
// It adds the keepGetCmd to the keepCmd and sets up the required flag "name".
// If the "name" flag is not provided, it logs an error using the slog package.
func init() {
    const op = "keep_get"

    keepCmd.AddCommand(keepGetCmd)

    keepGetCmd.Flags().String("name", "", "Secret name")
    if err := keepGetCmd.MarkFlagRequired("name"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }
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

    // Load the authentication token from the token storage.
    token, err := tokenStorage.Load()
    if err != nil {
        return
    }

    // Get the client configuration for connecting to the goph-keeper service.
    cfg := config.GetInstance().Config.Client

    // Create a new keeper client using the provided configuration and authentication token.
    keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

    // Request the secret from the goph-keeper service using the secret name.
    resp, err := keeperClient.GetItem(context.Background(), &v1.GetItemRequestV1{
        Name: name,
    })
    if err != nil {
        log.Error("Failed to create secret: ", slog.String("error", err.Error()))
    }

    // Decrypt the retrieved secret content.
    secret, err := decryptSecret(resp.GetContent())
    if err != nil {
        log.Error("Failed to decrypt secret: ", slog.String("error", err.Error()))
    }

    // Print the decrypted secret to the console.
    fmt.Printf("%s\n", secret)
}
