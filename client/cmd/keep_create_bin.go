/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// binCmd represents the bin command
var keepCreateBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep.create.bin"

		log := logger.GetInstance().Log.With("op", op)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error("Error reading secret name ", err)
			return
		}

		filePath, err := cmd.Flags().GetString("file_path")
		if err != nil {
			log.Error("Error reading file path ", err)
			return
		}

		token, err := tokenStorage.Load()
		if err != nil {
			return
		}

		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(token))
		err = keeperClient.CreateItemStream(context.Background(), name, filePath)
		if err != nil {
			log.Error("Error creating bin", err)
		}

	},
}

func init() {
	keepCreateCmd.AddCommand(keepCreateBinCmd)

	keepCreateBinCmd.Flags().String("name", "", "Secret name")
	if err := keepCreateBinCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ", err)
	}

	keepCreateBinCmd.Flags().StringP("file", "f", "", "Binary file path")
	if err := keepCreateBinCmd.MarkFlagRequired("file"); err != nil {
		slog.Error("Error setting flag: ", err)
	}
}
