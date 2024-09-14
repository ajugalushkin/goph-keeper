package register

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/auth"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// NewCommand creates a new cobra.Command for registering a user in the gophkeeper service.
// The command is configured with two flags: "email" and "password".
// The "email" flag is required and specifies the user's email address.
// The "password" flag is required and specifies the user's password.
// When the command is executed, the registerCmdRun function is called to handle the registration process.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Registers a user in the gophkeeper service",
		Run:   registerCmdRun,
	}

	cmd.Flags().StringP("email", "e", "", "User Email")
	if err := cmd.MarkFlagRequired("email"); err != nil {
		slog.Error("Error setting email flag", slog.String("error", err.Error()))
	}
	cmd.Flags().StringP("password", "p", "", "User password")
	if err := cmd.MarkFlagRequired("password"); err != nil {
		slog.Error("Error setting password flag", slog.String("error", err.Error()))
	}

	return cmd
}

// registerCmdRun handles the execution of the 'register' command.
// This function is responsible for registering a user in the gophkeeper service.
// It retrieves the user's email and password from the command-line flags,
// creates a new authentication client, and calls the Register method to register the user.
// If any errors occur during the process, they are logged and the function returns.
// Otherwise, a success message is printed.
func registerCmdRun(cmd *cobra.Command, args []string) {
	const op = "client.auth.register.run"
	log := logger.GetInstance().Log.With("op", op)

	// Retrieve the user's email from the command-line flag
	email, err := cmd.Flags().GetString("email")
	if err != nil {
		log.Error("Error getting email", slog.String("error", err.Error()))
		return
	}

	// Retrieve the user's password from the command-line flag
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		log.Error("Error getting password", slog.String("error", err.Error()))
		return
	}

	// Get the application configuration
	cfg := config.GetInstance().Config

	// Create a new authentication client using the application configuration
	authClient := auth.NewAuthClient(auth.GetAuthConnection(log, cfg.Client))

	// Register the user using the authentication client
	err = authClient.Register(context.Background(), email, password)
	if err != nil {
		log.Error("Error registering user", slog.String("error", err.Error()))
		return
	}

	// Print a success message
	fmt.Println("User registered successfully")
}
