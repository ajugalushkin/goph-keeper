package cmd

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// textCmd represents the text command
var keepUpdateTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Update text secret",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep update text"
		log := logger.GetInstance().Log.With("op", op)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error("Error reading secret name: ",
				slog.String("error", err.Error()))
			return
		}

		data, err := cmd.Flags().GetString("data")
		if err != nil {
			log.Error("Error reading text data: ",
				slog.String("error", err.Error()))
			return
		}

		text := vaulttypes.Text{
			Data: data,
		}

		content, err := encryptSecret(text)
		if err != nil {
			log.Error("Failed to encrypt secret: ",
				slog.String("error", err.Error()))
			return
		}

		token, err := tokenStorage.Load()
		if err != nil {
			return
		}

		cfg := config.GetInstance().Config.Client
		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

		resp, err := keeperClient.UpdateItem(context.Background(), &v1.UpdateItemRequestV1{
			Name:    name,
			Content: content,
		})
		if err != nil {
			log.Error("Failed to update secret: ",
				slog.String("error", err.Error()))
			return
		}

		fmt.Printf("Secret %s version %v updated successfully\n", resp.GetName(), resp.GetVersion())
	},
}

func init() {
	keepUpdateCmd.AddCommand(keepUpdateTextCmd)
}
