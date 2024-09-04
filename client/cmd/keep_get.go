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

// getCmd represents the get command
var keepGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secret",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep_get"
		log := logger.GetInstance().Log.With("op", op)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error("Error reading secret name: ", slog.String("error", err.Error()))
		}

		token, err := tokenStorage.Load()
		if err != nil {
			return
		}
		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(token))

		resp, err := keeperClient.GetItem(context.Background(), &v1.GetItemRequestV1{
			Name: name,
		})
		if err != nil {
			log.Error("Failed to create secret: ", slog.String("error", err.Error()))
		}

		secret, err := decryptVault(resp.GetContent())
		if err != nil {
			log.Error("Failed to decrypt secret: ", slog.String("error", err.Error()))
		}

		fmt.Printf("%s\n", secret)
	},
}

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
