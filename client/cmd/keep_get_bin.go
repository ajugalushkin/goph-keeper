package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// keepGetBinCmd represents the get command
var keepGetBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "Get secret",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep_get_bin"
		log := logger.GetInstance().Log.With("op", op)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error("Error reading secret name: ", slog.String("error", err.Error()))
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			log.Error("Error reading secret path: ", slog.String("error", err.Error()))
		}

		token, err := tokenStorage.Load()
		if err != nil {
			return
		}

		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(token))
		err = keeperClient.GetFile(context.Background(), name, path)
		if err != nil {
			log.Error("Failed to get secret: ", slog.String("error", err.Error()))
			return
		}

		fmt.Printf("file downloaded: %s\n", path)
	},
}

func init() {
	const op = "keep_get_bin"

	keepGetCmd.AddCommand(keepGetBinCmd)

	keepGetBinCmd.Flags().String("name", "", "Secret name")
	if err := keepGetBinCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	keepGetBinCmd.Flags().String("path", "", "Path to downloaded secret")
	if err := keepGetBinCmd.MarkFlagRequired("path"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
}
