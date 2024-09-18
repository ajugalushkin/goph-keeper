package register

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/auth"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

var registerCmd = NewCommand()

var client app.AuthClient

// NewCommand creates a new cobra.Command for registering a user in the gophkeeper service.
// It initializes a new Register instance with the provided logger and authentication client.
//
// Parameters:
// - newLog: A pointer to an slog.Logger instance used for logging messages.
// - newClient: An app.AuthClient instance used for authenticating users.
//
// Returns:
// - A pointer to the cobra.Command representing the 'register' command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Registers a user in the gophkeeper service",
		RunE:  registerCmdRun,
	}
}

// registerCmdRun handles the registration process for a user in the gophkeeper service.
// It retrieves the user's email and password from the command-line flags and registers the user using the provided authentication client.
//
// Parameters:
// - cmd: A pointer to the cobra.Command object representing the 'register' command.
// - args: A slice of strings containing any additional arguments passed to the command.
//
// Returns:
// - An error if any error occurs during the registration process. If no error occurs, it returns nil.
func registerCmdRun(cmd *cobra.Command, args []string) error {
	const op = "client.auth.register.run"
	log := logger.GetLogger().With("op", op)

	// Retrieve the user's email from the command-line flag
	email, err := cmd.Flags().GetString("email")
	if err != nil || email == "" {
		if err != nil {
			log.Error("Error getting email", slog.String("error", err.Error()))
		}
		return fmt.Errorf("email is required")
	}

	// Retrieve the user's password from the command-line flag
	password, err := cmd.Flags().GetString("password")
	if err != nil || password == "" {
		if err != nil {
			log.Error("Error getting password", slog.String("error", err.Error()))
		}
		return fmt.Errorf("password is required")
	}

	if client == nil {
		client = auth.NewAuthClient(auth.GetAuthConnection(log, config.GetConfig().Client))
	}

	// Register the user using the authentication client
	err = client.Register(context.Background(), email, password)
	if err != nil {
		log.Error("Error registering user", slog.String("error", err.Error()))
		return err
	}

	// Print a success message
	fmt.Println("User registered successfully")
	return nil
}

// registerCmdFlags configures the command-line flags for the 'register' command.
// It adds two flags to the provided cobra.Command: "email" and "password".
//
// The "email" flag is a string flag with a short name 'e' and a default value of "".
// It is used to specify the user's email address during registration.
//
// The "password" flag is a string flag with a short name 'p' and a default value of "".
// It is used to specify the user's password during registration.
//
// The function does not return any value.
func registerCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("email", "e", "", "User Email")
	cmd.Flags().StringP("password", "p", "", "User password")
}

func init() {
	registerCmdFlags(registerCmd)
}

func initClient(newClient app.AuthClient) {
	client = newClient
}
