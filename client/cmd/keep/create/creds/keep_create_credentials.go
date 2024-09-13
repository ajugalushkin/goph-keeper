package creds

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/secret"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// NewCommand creates a new Cobra command for creating creds secret.
// The command is configured to accept 'name', 'login', and 'password' flags.
// It also ensures that the 'name', 'login', and 'password' flags are required.
// Upon execution, the command calls the createCredentialsCmdRun function.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "creds",
		Short: "Create creds secret",
		Run:   createCredentialsCmdRun,
	}

	cmd.Flags().String("name", "", "Secret name")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Unable to mark 'name' flag as required %s", slog.String("error", err.Error()))
	}
	cmd.Flags().String("login", "", "Login")
	if err := cmd.MarkFlagRequired("login"); err != nil {
		slog.Error("Unable to mark 'login' flag as required %s", slog.String("error", err.Error()))
	}
	cmd.Flags().String("password", "", "Password")
	if err := cmd.MarkFlagRequired("password"); err != nil {
		slog.Error("Unable to mark 'password' flag as required %s", slog.String("error", err.Error()))
	}

	return cmd
}

// createCredentialsCmdRun handles the execution of the 'creds' command for creating a new secret.
// It retrieves the required parameters from the command-line flags, constructs a Credentials object,
// encrypts the object, and sends a request to the Keeper service to create the secret.
//
// Parameters:
// - cmd: The Cobra command object.
// - args: Additional command-line arguments.
//
// Returns:
// - No explicit return value. The function logs any errors encountered during execution and prints a success message upon completion.
func createCredentialsCmdRun(cmd *cobra.Command, args []string) {
	const op = "keep create creds"
	log := logger.GetInstance().Log.With("op", op)

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Unable to get `name` arg: ", slog.String("error", err.Error()))
		return
	}

	login, err := cmd.Flags().GetString("login")
	if err != nil {
		log.Error("Unable to get `login` arg: ", slog.String("error", err.Error()))
		return
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		log.Error("Unable to get `password` arg: ", slog.String("error", err.Error()))
		return
	}

	credentials := vaulttypes.Credentials{
		Login:    login,
		Password: password,
	}

	content, err := secret.EncryptSecret(credentials)
	if err != nil {
		log.Error("Failed to secret secret: ", slog.String("error", err.Error()))
		return
	}

	token, err := token_cache.GetInstance().Load()
	if err != nil {
		return
	}
	cfg := config.GetInstance().Config.Client
	keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

	resp, err := keeperClient.CreateItem(context.Background(), &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to create secret: ", slog.String("error", err.Error()))
		return
	}

	if resp == nil {
		log.Error("Failed to create secret: No response received")
		return
	}

	fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())
}
