package login

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"

	"github.com/spf13/cobra"
)

var (
	// log is used to log messages
	log *slog.Logger
	// login is used to login the user
	login *Login
)

type Login struct {
	client app.AuthClient
}

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
func NewCommand(newLog *slog.Logger, newClient app.AuthClient) *cobra.Command {
	log = newLog
	login = &Login{client: newClient}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Logins a user in the gophkeeper service",
		RunE:  loginCmdRun,
	}

	cmd.Flags().StringP("email", "e", "", "User Email")
	if err := cmd.MarkFlagRequired("email"); err != nil {
		slog.Error("Error marking email as required", slog.String("error", err.Error()))
	}
	cmd.Flags().StringP("password", "p", "", "User password")
	if err := cmd.MarkFlagRequired("password"); err != nil {
		slog.Error("Error marking password as required", slog.String("error", err.Error()))
	}

	return cmd
}

// loginCmdRun handles the execution of the login command.
// It retrieves the email and password from command-line flags, logs in the user using the provided credentials,
// saves the access token to the token cache, and prints it to the console.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
// - args: An array of strings representing command-line arguments.
//
// Return:
// - This function does not return any value.
func loginCmdRun(cmd *cobra.Command, args []string) error {
	const op = "client.auth.login.run"
	log.With("op", op)

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

	// Login the user using the provided email and password
	token, err := login.client.Login(context.Background(), email, password)
	if err != nil {
		log.Error("Error while login", slog.String("error", err.Error()))
		return err
	}

	// Save the access token to the token cache
	if err := token_cache.GetInstance().Save(token); err != nil {
		log.Error("Failed to store access token_cache", slog.String("error", err.Error()))
		return err
	}

	// Print the access token to the console
	fmt.Printf("Access Token: %s\n", token)
	return nil
}
