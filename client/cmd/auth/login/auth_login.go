package login

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/auth"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"

	"github.com/spf13/cobra"
)

// loginCmd is the command line for login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logins a user in the gophkeeper service",
	RunE:  authLoginCmdRunE,
}

var client app.AuthClient

// NewCommand creates a new Cobra command for logging in a user in the gophkeeper service.
// It takes a logger and an authentication client as parameters and returns a configured Cobra command.
// The command accepts email and password flags, logs in the user using the provided credentials,
// saves the access token to the token cache, and prints it to the console.
//
// Parameters:
// - newLog: A pointer to a slog.Logger object used for logging messages.
// - newClient: An implementation of the app.AuthClient interface for authenticating users.
//
// Return:
// - A pointer to a configured Cobra command for logging in a user.
func NewCommand() *cobra.Command {
	return loginCmd
}

// authLoginCmdRunE is the main function for the login command. It handles the user authentication process.
// It retrieves the email and password from command-line flags, logs in the user using the provided credentials,
// saves the access token to the token cache, and prints it to the console.
//
// Parameters:
// - cmd: A pointer to the Cobra command object. This object represents the login command and its associated flags.
// - args: A slice of strings containing any additional arguments passed to the command. In this case, it is not used.
//
// Return:
// - An error if any error occurs during the login process. If the login is successful, it returns nil.
func authLoginCmdRunE(cmd *cobra.Command, args []string) error {
	const op = "client.auth.login.run"
	log := logger.GetLogger().With("op", op)

	// Retrieve email from command-line flags
	email, err := cmd.Flags().GetString("email")
	if err != nil {
		log.Error("Error while getting email", slog.String("error", err.Error()))
		return err
	}
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Retrieve password from command-line flags
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		log.Error("Error while getting password", slog.String("error", err.Error()))
		return err
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}

	if client == nil {
		client = auth.NewAuthClient(auth.GetAuthConnection(log, config.GetConfig().Client))
	}

	// Login the user using the provided email and password
	token, err := client.Login(context.Background(), email, password)
	if err != nil {
		log.Error("Error while login", slog.String("error", err.Error()))
		return err
	}

	// Save the access token to the token cache
	if err := token_cache.GetToken().Save(token); err != nil {
		log.Error("Failed to store access token_cache", slog.String("error", err.Error()))
		return err
	}

	// Print the access token to the console
	fmt.Printf("Access Token: %s\n", token)
	return nil
}

// loginCmdFlags sets up the command-line flags for the login command.
// It adds two flags: "email" and "password". These flags are used to provide the user's email and password
// when running the login command.
//
// Parameters:
// - cmd: A pointer to the Cobra command object. This object represents the login command and its associated flags.
//
// Return:
// - This function does not return any value. It modifies the provided Cobra command object directly.
func loginCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("email", "e", "", "User Email")
	cmd.Flags().StringP("password", "p", "", "User password")
}

func init() {
	loginCmdFlags(loginCmd)
}

func initClient(newClient app.AuthClient) {
	client = newClient
}
