package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
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
			return
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			log.Error("Error reading secret path: ", slog.String("error", err.Error()))
			return
		}

		token, err := tokenStorage.Load()
		if err != nil {
			log.Error("Error loading token: ", slog.String("error", err.Error()))
			return
		}

		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(token))
		stream, err := keeperClient.GetFile(context.Background(), name)
		if err != nil {
			log.Error("Error getting file stream: ", slog.String("error", err.Error()))
			return
		}

		req, err := stream.Recv()
		if err != nil {
			log.Error("Error getting file info: ", slog.String("error", err.Error()))
			return
		}

		secret, err := decryptSecret(req.GetContent())
		if err != nil {
			log.Error("Failed to decrypt secret: ", slog.String("error", err.Error()))
		}

		fileInfo := secret.(vaulttypes.Bin)

		newFile, err := os.Create(filepath.Join(path, fileInfo.FileName))
		if err != nil {
			log.Error("Error creating file: ", slog.String("error", err.Error()))
		}
		defer newFile.Close()

		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Error("Error getting file chunk: ", slog.String("error", err.Error()))
				return
			}
			chunk := req.GetChunkData()

			_, err = newFile.Write(chunk)
			if err != nil {
				log.Error("Error add chunk to file: ", slog.String("error", err.Error()))
				return
			}
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
