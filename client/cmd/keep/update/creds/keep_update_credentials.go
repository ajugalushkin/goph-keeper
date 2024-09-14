package creds

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/keeper"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/secret"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// NewCommand creates a new Cobra command for updating creds secret in the Goph-Keeper vault.
// The command accepts three flags: --name, --login, and --password.
// It validates the flags and calls the keepUpdateCredRun function to perform the update operation.
func NewCommand() *cobra.Command {
	const op = "keep update creds"

	cmd := &cobra.Command{
		Use:   "creds",
		Short: "Update creds secret",
		Run:   keepUpdateCredRun,
	}

	// Flag for specifying the secret name.
	cmd.Flags().String("name", "", "Secret name")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	// Flag for specifying the login.
	cmd.Flags().String("login", "", "Login")
	if err := cmd.MarkFlagRequired("login"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	// Flag for specifying the password.
	cmd.Flags().String("password", "", "Password")
	if err := cmd.MarkFlagRequired("password"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	return cmd
}

// keepUpdateCredRun is a function that updates a secret in the Goph-Keeper vault.
// It reads the secret name, login, and password from command-line flags, encrypts the creds,
// and sends an update request to the Goph-Keeper server.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
// - args: An array of command-line arguments.
//
// Return:
// - This function does not return any value.
func keepUpdateCredRun(cmd *cobra.Command, args []string) {
	const op = "keep update creds"
	log := logger.GetLogger().With("op", op)

	// Read the secret name from the command-line flag.
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ",
			slog.String("error", err.Error()))
	}

	// Read the login from the command-line flag.
	login, err := cmd.Flags().GetString("login")
	if err != nil {
		log.Error("Error reading login: ",
			slog.String("error", err.Error()))
	}

	// Read the password from the command-line flag.
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		log.Error("Error reading password: ",
			slog.String("error", err.Error()))
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
		return
	}

	// Load the authentication token_cache from the token_cache storage.
	token, err := token_cache.GetToken().Load()
	if err != nil {
		return
	}

	// Get the Goph-Keeper client configuration.
	cfg := config.GetConfig().Client

	// Create a new Goph-Keeper client using the provided configuration and authentication token_cache.
	keeperClient := keeper.NewKeeperClient(keeper.GetKeeperConnection(log, cfg.Address, token))

	// Send an update request to the Goph-Keeper server with the secret name and encrypted content.
	resp, err := keeperClient.UpdateItem(context.Background(), &v1.UpdateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to update secret: ",
			slog.String("error", err.Error()))
		return
	}

	// Print a success message with the updated secret name and version.
	fmt.Printf("Secret %s version %v updated successfully\n", resp.GetName(), resp.GetVersion())
}
