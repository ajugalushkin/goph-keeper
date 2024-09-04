package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// keepCreateTextCmd represents the text command
var keepCreateTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Create text secret",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep_create_text"
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

		content, err := encryptVault(text)
		if err != nil {
			log.Error("Failed to encrypt secret: ",
				slog.String("error", err.Error()))
			return
		}

		token, err := tokenStorage.Load()
		if err != nil {
			return
		}
		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(token))

		resp, err := keeperClient.CreateItem(context.Background(), &v1.CreateItemRequestV1{
			Name:    name,
			Content: content,
		})
		if err != nil {
			log.Error("Failed to create secret: ", slog.String("error", err.Error()))
		}

		fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())
	},
}

func init() {
	const op = "keep_create_text"
	keepCreateCmd.AddCommand(keepCreateTextCmd)

	keepCreateTextCmd.Flags().String("name", "", "Secret name")
	if err := keepCreateTextCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepCreateTextCmd.Flags().String("data", "", "Text data")
	if err := keepCreateTextCmd.MarkFlagRequired("data"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
}
