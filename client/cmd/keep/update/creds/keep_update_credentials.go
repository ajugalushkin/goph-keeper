package creds

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/ajugalushkin/goph-keeper/client/secret"
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

var keepUpdateCreds = &cobra.Command{
	Use:   "creds",
	Short: "Update creds secret",
	RunE:  keepUpdateCredsRunE,
}

var client app.KeeperClient

// NewCommand creates a new Cobra command for updating creds secret in the Goph-Keeper vault.
// The command accepts three flags: --name, --login, and --password.
// It validates the flags and calls the keepUpdateCredRun function to perform the update operation.
func NewCommand() *cobra.Command {
	return keepUpdateCreds
}

// keepUpdateCredsRunE is the entry point for updating a secret in the Goph-Keeper vault.
// It reads the secret name, login, and password from command-line flags, creates a Credentials object,
// encrypts the credentials, loads the authentication token, creates a Goph-Keeper client, sends an update request,
// and prints a success message upon successful update.
//
// Parameters:
// - cmd: A pointer to the Cobra command object. This object represents the command and its associated flags.
// - args: A slice of strings representing any additional arguments passed to the command.
//
// Return:
//   - An error if any error occurs during the execution of the command.
//     If no error occurs, it returns nil.
func keepUpdateCredsRunE(cmd *cobra.Command, args []string) error {
	const op = "keep update creds"
	log := logger.GetLogger().With("op", op)

	// Read the secret name from the command-line flag.
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ",
			slog.String("error", err.Error()))
		return fmt.Errorf("error reading secret name")
	}
	if name == "" {
		return fmt.Errorf("name is required")
	}

	// Read the login from the command-line flag.
	login, err := cmd.Flags().GetString("login")
	if err != nil {
		log.Error("Error reading login: ",
			slog.String("error", err.Error()))
		return fmt.Errorf("error reading login")
	}
	if login == "" {
		return fmt.Errorf("login is required")
	}

	// Read the password from the command-line flag.
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		log.Error("Error reading password: ",
			slog.String("error", err.Error()))
		return fmt.Errorf("error reading password")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}

	// Create a Credentials object with the provided login and password.
	credentials := vaulttypes.Credentials{
		Login:    login,
		Password: password,
	}

	// Encrypt the creds using the encryptSecret function.
	content, err := secret.EncryptSecret(credentials)
	if err != nil {
		log.Error("Failed to secret secret: ",
			slog.String("error", err.Error()))
		return err
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
	// Send an update request to the Goph-Keeper server with the secret name and encrypted content.
	resp, err := client.UpdateItem(context.Background(), &v1.UpdateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to update secret: ",
			slog.String("error", err.Error()))
		return err
	}

	// Print a success message with the updated secret name and version.
	fmt.Printf("Secret %s version %v updated successfully\n", resp.GetName(), resp.GetVersion())
	return nil
}

// updateCredsCmdFlags sets up command-line flags for updating a secret in the Goph-Keeper vault.
// The function accepts a pointer to a Cobra command object and adds three flags: --name, --login, and --password.
// These flags are used to provide the secret name, login, and password for updating the secret in the vault.
//
// Parameters:
// - cmd: A pointer to the Cobra command object. This object represents the command and its associated flags.
//
// Return:
// - This function does not return any value. It modifies the provided Cobra command object by adding the required flags.
func updateCredsCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Secret name")
	cmd.Flags().String("login", "", "Login")
	cmd.Flags().String("password", "", "Password")
}

func init() {
	updateCredsCmdFlags(keepUpdateCreds)
}

func initClient(newClient app.KeeperClient) {
	client = newClient
}
