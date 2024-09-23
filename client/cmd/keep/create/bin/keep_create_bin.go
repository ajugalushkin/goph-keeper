package bin

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
)

var keepCreateBin = &cobra.Command{
	Use:   "bin",
	Short: "Create bin secret",
	RunE:  keepCreateBinCmdRunE,
}

var client app.KeeperClient

// NewCommand creates a new Cobra command for creating a bin secret.
// It initializes the logger and client for the command and returns the command object.
//
// Parameters:
// - newLog: A pointer to an slog.Logger object used for logging messages.
// - newClient: An implementation of the app.KeeperClient interface for interacting with the vault.
//
// Return:
// - A pointer to the Cobra command object for creating a bin secret.
func NewCommand() *cobra.Command {
	return keepCreateBin
}

// keepCreateBinCmdRunE is the entry point for the "bin" command in the "keep create" command group.
// It reads the secret name and file path from the command-line flags, prepares the file information,
// encrypts the secret content, opens the file for reading, creates the binary secret in the vault,
// and prints a success message.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
// - args: An array of strings representing command-line arguments.
//
// Return:
//   - An error if any error occurs during the execution of the command.
//     If no error occurs, it returns nil.
func keepCreateBinCmdRunE(cmd *cobra.Command, args []string) error {
	const op = "keep.create.bin"
	log := logger.GetLogger().With("op", op)

	// Read the secret name from the command-line flags
	name, err := cmd.Flags().GetString("name")
	if err != nil || len(name) == 0 {
		if err != nil {
			log.Error("Error reading secret name ", slog.String("error", err.Error()))
		}
		return fmt.Errorf("name is required")
	}

	// Read the file path from the command-line flags
	filePath, err := cmd.Flags().GetString("file_path")
	if err != nil || len(filePath) == 0 {
		if err != nil {
			log.Error("Error reading file path ", slog.String("error", err.Error()))
		}
		return fmt.Errorf("file_path is required")
	}

	if client == nil {
		token, err := token_cache.GetToken().Load()
		if err != nil {
			return err
		}
		client = keeper.NewKeeperClient(keeper.GetKeeperConnection(log, config.GetConfig().Client.Address, token))
	}

	// Create the binary secret in the vault
	resp, err := client.CreateItemStream(context.Background(), name, filePath)
	if err != nil {
		log.Error("Error creating bin", slog.String("error", err.Error()))
		return err
	}

	// Print the success message
	fmt.Printf("Secret %s created successfully\n", resp.GetName())
	return nil
}

// binCmdFlags sets up the flags for the "bin" command in the "keep create" command group.
// The flags are used to specify the name and file path of the binary secret to be created.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
//
// Return:
// - None.
func binCmdFlags(cmd *cobra.Command) {
	// Add a flag for the secret name. The default value is an empty string.
	cmd.Flags().StringP("name", "n", "", "Secret name")

	// Add a flag for the file path. The default value is an empty string.
	// The flag is also available with a short name "f".
	cmd.Flags().StringP("file_path", "f", "", "Binary file path")
}

func init() {
	binCmdFlags(keepCreateBin)
}

func initClient(newClient app.KeeperClient) {
	client = newClient
}
