package cmd

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
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
	Run:   keepCreateBinCmdRun,
}

// init initializes the keepCreateBinCmd command and its flags.
// It adds the keepCreateBinCmd command to the keepCreateCmd command and sets up the required flags.
// The flags include:
// - "name": The name of the secret to be created. This flag is required.
// - "file_path" (or "f"): The path to the binary file that will be stored as a secret. This flag is required.
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

// keepCreateBinCmdRun is the main function for the "bin" command in the "keep create" command group.
// It creates a new binary secret in the Goph-Keeper vault.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
// - args: An array of strings representing the command-line arguments.
//
// Return:
// - None.
func keepCreateBinCmdRun(cmd *cobra.Command, args []string) {
	const op = "keep.create.bin"
	log := logger.GetInstance().Log.With("op", op)

	// Read the secret name from the command-line flags
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name ", slog.String("error", err.Error()))
		return
	}

	// Read the file path from the command-line flags
	filePath, err := cmd.Flags().GetString("file_path")
	if err != nil {
		log.Error("Error reading file path ", slog.String("error", err.Error()))
		return
	}

	// Get the file statistics
	stat, err := os.Stat(filePath)
	if err != nil {
		log.Error("Error reading file stat ", slog.String("error", err.Error()))
		return
	}

	// Prepare the file information for the secret
	fileInfo := vaulttypes.Bin{
		FileName: filepath.Base(filePath),
		Size:     stat.Size(),
	}

	// Encrypt the secret content
	content, err := encryptSecret(fileInfo)
	if err != nil {
		log.Error("Failed to encrypt secret: ", slog.String("error", err.Error()))
		return
	}

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		log.Error("cannot open file: ", slog.String("error", err.Error()))
		return
	}
	defer file.Close()

	// Load the authentication token
	token, err := tokenStorage.Load()
	if err != nil {
		log.Error("Error loading token: ", slog.String("error", err.Error()))
		return
	}

	// Get the Goph-Keeper client configuration
	cfg := config.GetInstance().Config.Client

	// Create a new Goph-Keeper client
	keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

	// Create the binary secret in the vault
	resp, err := keeperClient.CreateItemStream(context.Background(), name, file, content)
	if err != nil {
		log.Error("Error creating bin", slog.String("error", err.Error()))
		return
	}

	// Print the success message
	fmt.Printf("Secret %s created successfully\n", resp.GetName())
}
