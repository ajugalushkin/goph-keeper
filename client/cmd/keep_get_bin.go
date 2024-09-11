package cmd

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
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
	Run: keepGetBinRun,
}

// init initializes the keepGetBinCmd command and its flags.
// It adds the command to the keepGetCmd command and sets up required flags.
// The flags include:
// - "name": The name of the secret to be retrieved.
// - "path": The path where the downloaded secret will be saved.
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
    log := logger.GetInstance().Log.With("op", op)

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

    // Load the authentication token from the token storage
    token, err := tokenStorage.Load()
    if err != nil {
        log.Error("Error loading token: ", slog.String("error", err.Error()))
        return
    }

    // Create a new gRPC client for interacting with the goph-keeper service
    cfg := config.GetInstance().Config.Client
    keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

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
    secret, err := decryptSecret(req.GetContent())
    if err != nil {
        log.Error("Failed to decrypt secret: ", slog.String("error", err.Error()))
    }

    // Extract the file information from the decrypted secret
    fileInfo := secret.(vaulttypes.Bin)

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