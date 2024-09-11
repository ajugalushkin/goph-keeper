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

// credentialsCmd represents the credentials command
var keepUpdateCredentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Update credentials secret",
	Run: keepUpdateCredRun,
}

// init initializes the "credentials" command for the "keep update" command.
// It sets up the command flags for specifying the secret name, login, and password.
// If any required flag is missing, it logs an error message.
func init() {
    const op = "keep update credentials"
    keepUpdateCmd.AddCommand(keepUpdateCredentialsCmd)

    // Flag for specifying the secret name.
    keepUpdateCredentialsCmd.Flags().String("name", "", "Secret name")
    if err := keepUpdateCredentialsCmd.MarkFlagRequired("name"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }

    // Flag for specifying the login.
    keepUpdateCredentialsCmd.Flags().String("login", "", "Login")
    if err := keepUpdateCredentialsCmd.MarkFlagRequired("login"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }

    // Flag for specifying the password.
    keepUpdateCredentialsCmd.Flags().String("password", "", "Password")
    if err := keepUpdateCredentialsCmd.MarkFlagRequired("password"); err != nil {
        slog.Error("Error setting flag: ",
            slog.String("op", op),
            slog.String("error", err.Error()))
    }
}

// keepUpdateCredRun is a function that updates a secret in the Goph-Keeper vault.
// It reads the secret name, login, and password from command-line flags, encrypts the credentials,
// and sends an update request to the Goph-Keeper server.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
// - args: An array of command-line arguments.
//
// Return:
// - This function does not return any value.
func keepUpdateCredRun(cmd *cobra.Command, args []string) {
    const op = "keep update credentials"
    log := logger.GetInstance().Log.With("op", op)

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

    // Encrypt the credentials using the encryptSecret function.
    content, err := encryptSecret(credentials)
    if err != nil {
        log.Error("Failed to encrypt secret: ",
            slog.String("error", err.Error()))
        return
    }

    // Load the authentication token from the token storage.
    token, err := tokenStorage.Load()
    if err != nil {
        return
    }

    // Get the Goph-Keeper client configuration.
    cfg := config.GetInstance().Config.Client

    // Create a new Goph-Keeper client using the provided configuration and authentication token.
    keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

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