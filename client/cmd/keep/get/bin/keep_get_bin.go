package bin

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/secret"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
)

// NewCommand creates a new cobra.Command for retrieving a secret file from the goph-keeper service.
// The command is configured to accept two flags: "name" and "path".
// The "name" flag specifies the name of the secret to be retrieved.
// The "path" flag specifies the path where the downloaded secret will be saved.
// If either flag is not provided, an error will be logged.
// The command runs the keepGetBinRun function to handle the retrieval process.
func NewCommand() *cobra.Command {
	const op = "keep_get_bin"

	cmd := &cobra.Command{
		Use:   "bin",
		Short: "Get secret",
		Run:   keepGetBinRun,
	}

	cmd.Flags().String("name", "", "Secret name")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	cmd.Flags().String("path", "", "Path to downloaded secret")
	if err := cmd.MarkFlagRequired("path"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	return cmd
}

// keepGetBinRun is the function that handles the retrieval of a secret file from the goph-keeper service.
// It uses gRPC to communicate with the server and streams the file data to a local file.
//
// Parameters:
// - cmd: A pointer to the cobra.Command object representing the command being executed.
// - args: A slice of strings containing any additional arguments passed to the command.
//
// Return:
// - This function does not return any value.
func keepGetBinRun(cmd *cobra.Command, args []string) {
	const op = "keep_get_bin"
	log := logger.GetLogger().With("op", op)

	// Read the secret name from the command flags
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ", slog.String("error", err.Error()))
		return
	}

	// Read the path where the downloaded secret will be saved from the command flags
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		log.Error("Error reading secret path: ", slog.String("error", err.Error()))
		return
	}

	// Load the authentication token_cache from the token_cache storage
	token, err := token_cache.GetToken().Load()
	if err != nil {
		log.Error("Error loading token_cache: ", slog.String("error", err.Error()))
		return
	}

	// Create a new gRPC client for interacting with the goph-keeper service
	cfg := config.GetConfig().Client
	keeperClient := keeper.NewKeeperClient(keeper.GetKeeperConnection(log, cfg.Address, token))

	// Request the file stream from the goph-keeper service
	stream, err := keeperClient.GetFile(context.Background(), name)
	if err != nil {
		log.Error("Error getting file stream: ", slog.String("error", err.Error()))
		return
	}

	// Receive the file information from the stream
	req, err := stream.Recv()
	if err != nil {
		log.Error("Error getting file info: ", slog.String("error", err.Error()))
		return
	}

	// Decrypt the secret file content
	respSecret, err := secret.DecryptSecret(req.GetContent())
	if err != nil {
		log.Error("Failed to decrypt secret: ", slog.String("error", err.Error()))
	}

	// Extract the file information from the decrypted secret
	fileInfo := respSecret.(vaulttypes.Bin)

	// Create a new local file to save the downloaded secret
	newFile, err := os.Create(filepath.Join(path, fileInfo.FileName))
	if err != nil {
		log.Error("Error creating file: ", slog.String("error", err.Error()))
	}
	defer newFile.Close()

	// Stream the file chunks from the goph-keeper service to the local file
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

	// Print a success message
	fmt.Printf("file downloaded: %s\n", path)
}
