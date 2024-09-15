package bin

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/secret"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
)

var keepGetBin = &cobra.Command{
	Use:   "bin",
	Short: "Get secret",
	RunE:  keepGetBinRunE,
}

var client app.KeeperClient

// NewCommand creates a new cobra.Command for retrieving a secret file from the goph-keeper service.
// The command is configured to accept two flags: "name" and "path".
// The "name" flag specifies the name of the secret to be retrieved.
// The "path" flag specifies the path where the downloaded secret will be saved.
// If either flag is not provided, an error will be logged.
// The command runs the keepGetBinRun function to handle the retrieval process.
func NewCommand() *cobra.Command {
	return keepGetBin
}

// keepGetBinRunE is a function that handles the retrieval of a secret file from the goph-keeper service.
// It reads the secret name and path from command flags, initializes the Keeper client if not already done,
// requests the file stream from the goph-keeper service, decrypts the secret file content,
// and streams the file chunks to a local file.
//
// Parameters:
// - cmd: A pointer to the Cobra command object representing the "get bin" command.
// - args: An array of strings representing any additional arguments passed to the command.
//
// Returns:
// - An error if any error occurs during the retrieval process.
// - nil if the retrieval is successful.
func keepGetBinRunE(cmd *cobra.Command, args []string) error {
	const op = "keep_get_bin"
	log := logger.GetLogger().With("op", op)

	// Read the secret name from the command flags
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ", slog.String("error", err.Error()))
		return fmt.Errorf("error reading secret name")
	}
	if name == "" {
		return fmt.Errorf("secret name is required")
	}

	// Read the path where the downloaded secret will be saved from the command flags
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		log.Error("Error reading secret path: ", slog.String("error", err.Error()))
		return fmt.Errorf("error reading secret path")
	}
	if path == "" {
		return fmt.Errorf("secret path is required")
	}

	// If the Keeper client is not initialized, load the authentication token_cache from storage and create a new client.
	if client == nil {
		// Load the authentication token_cache from storage.
		token, err := token_cache.GetToken().Load()
		if err != nil {
			return err
		}
		// Create a new Keeper client using the provided configuration and token_cache.
		cfg := config.GetConfig().Client
		client = keeper.NewKeeperClient(keeper.GetKeeperConnection(log, cfg.Address, token))
	}
	// Request the file stream from the goph-keeper service
	stream, err := client.GetFile(context.Background(), name)
	if err != nil {
		log.Error("Error getting file stream: ", slog.String("error", err.Error()))
		return err
	}

	// Receive the file information from the stream
	req, err := stream.Recv()
	if err != nil {
		log.Error("Error getting file info: ", slog.String("error", err.Error()))
		return err
	}

	// Decrypt the secret file content
	respSecret, err := secret.DecryptSecret(req.GetContent())
	if err != nil {
		log.Error("Failed to decrypt secret: ", slog.String("error", err.Error()))
		return err
	}

	// Extract the file information from the decrypted secret
	fileInfo := respSecret.(vaulttypes.Bin)

	// Create a new local file to save the downloaded secret
	newFile, err := os.Create(filepath.Join(path, fileInfo.FileName))
	if err != nil {
		log.Error("Error creating file: ", slog.String("error", err.Error()))
		return err
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
			return err
		}
		chunk := req.GetChunkData()

		_, err = newFile.Write(chunk)
		if err != nil {
			log.Error("Error add chunk to file: ", slog.String("error", err.Error()))
			return err
		}
	}

	// Print a success message
	fmt.Printf("file downloaded: %s\n", path)
	return nil
}

// getBinCmdFlags sets up the command line flags for the "get bin" command.
// It accepts two flags: "name" and "path".
//
// The "name" flag specifies the name of the secret to be retrieved.
// The "path" flag specifies the path where the downloaded secret will be saved.
//
// Parameters:
// - cmd: A pointer to the Cobra command object representing the "get bin" command.
//
// Returns:
// - None.
func getBinCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Secret name")
	cmd.Flags().String("path", "", "Path to downloaded secret")
}

func init() {
	getBinCmdFlags(keepGetBin)
}

func initClient(newClient app.KeeperClient) {
	client = newClient
}
