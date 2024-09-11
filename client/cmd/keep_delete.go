package cmd

import (
	"context"
	"fmt"

	"github.com/ajugalushkin/goph-keeper/client/config"
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
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep_delete"
		log := logger.GetInstance().Log.With("op", op)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error("Error reading secret name", slog.String("error", err.Error()))
		}

		token, err := tokenStorage.Load()
		if err != nil {
			return
		}

		cfg := config.GetInstance().Config.Client
		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

		resp, err := keeperClient.DeleteItem(context.Background(), &keeperv1.DeleteItemRequestV1{
			Name: name,
		})
		if err != nil {
			log.Error("Failed to delete secret", slog.String("error", err.Error()))
			return
		}

		fmt.Printf("Secret %s deleted successfully\n", resp.GetName())
	},
}

func init() {
	const op = "keep_create_card"
	keepCmd.AddCommand(keepDeleteCmd)

	keepDeleteCmd.Flags().String("name", "", "Secret name")
	if err := keepDeleteCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
}
