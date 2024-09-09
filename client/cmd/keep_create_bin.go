package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
)

// binCmd represents the bin command
var keepCreateBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "Create bin secret",
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

		stat, err := os.Stat(filePath)
		if err != nil {
			log.Error("Error reading file stat ", err)
			return
		}

		fileInfo := vaulttypes.Bin{
			FileName: filepath.Base(filePath),
			Size:     stat.Size(),
		}

		content, err := encryptSecret(fileInfo)
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
		resp, err := keeperClient.CreateItemStream(context.Background(), name, filePath, content)
		if err != nil {
			log.Error("Error creating bin", err)
			return
		}

		fmt.Printf("Secret %s version %v created successfully\n", resp.GetName())
	},
}

func init() {
	keepCreateCmd.AddCommand(keepCreateBinCmd)

	keepCreateBinCmd.Flags().String("name", "", "Secret name")
	if err := keepCreateBinCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ", slog.String("error", err.Error()))
	}

	keepCreateBinCmd.Flags().StringP("file_path", "f", "", "Binary file path")
	if err := keepCreateBinCmd.MarkFlagRequired("file_path"); err != nil {
		slog.Error("Error setting flag: ", slog.String("error", err.Error()))
	}
}
