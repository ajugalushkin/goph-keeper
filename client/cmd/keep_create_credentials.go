package cmd

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// keepCreateCredentialsCmd represents the credentials command
var keepCreateCredentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Create credentials secret",
	Run:   createCredentialsCmdRun,
}

// init initializes the 'credentials' command for the 'keep create' command.
// It adds the 'credentials' command to the 'keep create' command and sets up the required flags.
// The flags include 'name', 'login', and 'password'. The 'name' flag is required and specifies the name of the secret.
// The 'login' and 'password' flags are also required and represent the login and password for the secret, respectively.
// If any of the required flags are not provided, an error message is logged.
func init() {
	keepCreateCmd.AddCommand(keepCreateCredentialsCmd)

	keepCreateCredentialsCmd.Flags().String("name", "", "Secret name")
	if err := keepCreateCredentialsCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Unable to mark 'name' flag as required %s", slog.String("error", err.Error()))
	}
	keepCreateCredentialsCmd.Flags().String("login", "", "Login")
	if err := keepCreateCredentialsCmd.MarkFlagRequired("login"); err != nil {
		slog.Error("Unable to mark 'login' flag as required %s", slog.String("error", err.Error()))
	}
	keepCreateCredentialsCmd.Flags().String("password", "", "Password")
	if err := keepCreateCredentialsCmd.MarkFlagRequired("password"); err != nil {
		slog.Error("Unable to mark 'password' flag as required %s", slog.String("error", err.Error()))
	}
}

// createCredentialsCmdRun handles the execution of the 'credentials' command for creating a new secret.
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
	const op = "keep create credentials"
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

	content, err := encryptSecret(credentials)
	if err != nil {
		log.Error("Failed to encrypt secret: ", slog.String("error", err.Error()))
		return
	}

	token, err := tokenStorage.Load()
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
	}

	fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())
}
