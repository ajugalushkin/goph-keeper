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

var keepCreateCreds = &cobra.Command{
	Use:   "creds",
	Short: "Create creds secret",
	RunE:  createCredentialsCmdRunE,
}

var client app.KeeperClient

// NewCommand creates a new Cobra command for creating creds secret.
// The command is configured to accept 'name', 'login', and 'password' flags.
// It also ensures that the 'name', 'login', and 'password' flags are required.
// Upon execution, the command calls the createCredentialsCmdRun function.
func NewCommand() *cobra.Command {
	return keepCreateCreds
}

// createCredentialsCmdRunE is the entry point for creating a new credentials secret in the Goph-Keeper vault.
// It accepts command-line arguments and flags, validates the input, encrypts the credentials, and creates the secret using the Keeper client.
//
// Parameters:
// - cmd: A pointer to the Cobra command object that triggered the execution of this function.
// - args: An array of strings containing any additional arguments passed to the command.
//
// Returns:
// - An error if any step in the process fails, or nil if the secret is created successfully.
func createCredentialsCmdRunE(cmd *cobra.Command, args []string) error {
	const op = "keep create creds"
	log := logger.GetLogger().With("op", op)

	// Retrieve the 'name', 'login', and 'password' flags from the command-line arguments.
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Unable to get `name` arg: ", slog.String("error", err.Error()))
		return err
	}

	login, err := cmd.Flags().GetString("login")
	if err != nil {
		log.Error("Unable to get `login` arg: ", slog.String("error", err.Error()))
		return err
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		log.Error("Unable to get `password` arg: ", slog.String("error", err.Error()))
		return err
	}

	// Create a new Credentials object with the provided login and password.
	credentials := vaulttypes.Credentials{
		Login:    login,
		Password: password,
	}

	// Encrypt the credentials using the secret.EncryptSecret function.
	content, err := secret.EncryptSecret(credentials)
	if err != nil {
		log.Error("Failed to secret secret: ", slog.String("error", err.Error()))
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

	// Use the Keeper client to create a new item (secret) with the provided name and encrypted content.
	resp, err := client.CreateItem(context.Background(), &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to create secret: ", slog.String("error", err.Error()))
		return err
	}

	// Check if the response from the Keeper client is nil.
	if resp == nil {
		log.Error("Failed to create secret: No response received")
		return err
	}

	// Print a success message with the name and version of the newly created secret.
	fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())
	return nil
}

// credsCmdFlags configures the command-line flags for the 'text' command.
// It adds three flags: 'name', 'login', and 'password'.
//
// Parameters:
// - cmd: A pointer to the Cobra command object to which the flags will be added.
//
// The 'name' flag is used to specify the name of the secret to be created.
// The 'login' flag is used to store the login information for the secret.
// The 'password' flag is used to store the password information for the secret.
//
// The flags are configured with default values of "", and their descriptions are provided.
func credsCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Secret name")
	cmd.Flags().String("login", "", "Login")
	cmd.Flags().String("password", "", "Password")
}

func init() {
	credsCmdFlags(keepCreateCreds)
}

func initClient(newClient app.KeeperClient) {
	client = newClient
}
