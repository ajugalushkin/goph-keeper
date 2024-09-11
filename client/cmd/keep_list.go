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

// keepListCmd represents the list command
var keepListCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets",
	Run: keepListRun,
}

func init() {
	keepCmd.AddCommand(keepListCmd)
}

// keepListRun is a command handler function that lists secrets stored in the goph-keeper service.
// It retrieves the access token, establishes a connection to the goph-keeper service, sends a request to list secrets,
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

    token, err := tokenStorage.Load()
    if err != nil {
        return
    }

    cfg := config.GetInstance().Config.Client
    keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))
    resp, err := keeperClient.ListItems(context.Background(), &v1.ListItemsRequestV1{})
    if err != nil {
        log.Error("Failed to list secret: ", slog.String("error", err.Error()))
    }

    for _, info := range resp.GetSecrets() {
        secret, err := decryptSecret(info.GetContent())
        if err != nil {
            log.Error("Failed to decrypt secret: ", slog.String("error", err.Error()))
        }

        fmt.Printf("%s\n", secret)
    }
}
