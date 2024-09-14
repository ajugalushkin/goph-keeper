package bin

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/secret"

	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
)

var (
	// keepCreateBin is a command to create a bin data
	keepCreateBin = &cobra.Command{
		Use:   "bin",
		Short: "Create bin secret",
		RunE:  keepCreateBinCmdRunE,
	}

	// log is used to log messages
	log *slog.Logger

	bin *Bin
)

// Bin is used to set client
type Bin struct {
	client app.KeeperClient
}


// NewCommand creates a new Cobra command for creating a bin secret.
// It initializes the logger and client for the command and returns the command object.
//
// Parameters:
// - newLog: A pointer to an slog.Logger object used for logging messages.
// - newClient: An implementation of the app.KeeperClient interface for interacting with the vault.
//
// Return:
// - A pointer to the Cobra command object for creating a bin secret.
func NewCommand(newLog *slog.Logger, newClient app.KeeperClient) *cobra.Command {
    log = newLog
    bin = &Bin{client: newClient}

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
// - An error if any error occurs during the execution of the command.
//   If no error occurs, it returns nil.
func keepCreateBinCmdRunE(cmd *cobra.Command, args []string) error {
    const op = "keep.create.bin"
    log.With("op", op)

    // Read the secret name from the command-line flags
    name, err := cmd.Flags().GetString("name")
    if err != nil {
        log.Error("Error reading secret name ", slog.String("error", err.Error()))
        return err
    }

    // Read the file path from the command-line flags
    filePath, err := cmd.Flags().GetString("file_path")
    if err != nil {
        log.Error("Error reading file path ", slog.String("error", err.Error()))
        return err
    }

    // Get the file statistics
    stat, err := os.Stat(filePath)
    if err != nil {
        log.Error("Error reading file stat ", slog.String("error", err.Error()))
        return err
    }

    // Prepare the file information for the secret
    fileInfo := vaulttypes.Bin{
        FileName: filepath.Base(filePath),
        Size:     stat.Size(),
    }

    // Encrypt the secret content
    content, err := secret.EncryptSecret(fileInfo)
    if err != nil {
        log.Error("Failed to secret secret: ", slog.String("error", err.Error()))
        return err
    }

    // Open the file for reading
    file, err := os.Open(filePath)
    if err != nil {
        log.Error("cannot open file: ", slog.String("error", err.Error()))
        return err
    }
    defer file.Close()

    // Create the binary secret in the vault
    resp, err := bin.client.CreateItemStream(context.Background(), name, file, content)
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
    cmd.Flags().String("name", "", "Secret name")

    // Add a flag for the file path. The default value is an empty string.
    // The flag is also available with a short name "f".
    cmd.Flags().StringP("file_path", "f", "", "Binary file path")
}

func init() {
	binCmdFlags(keepCreateBin)
}
