package cmd

import (
	"context"
	"fmt"
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
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep_get"
		log := logger.GetInstance().Log.With("op", op)

		token, err := tokenStorage.Load()
		if err != nil {
			return
		}

		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(token))
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
	},
}

func init() {
	keepCmd.AddCommand(keepListCmd)
}
