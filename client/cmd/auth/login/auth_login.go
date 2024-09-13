package login

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// NewCommand creates a new Cobra command for user login to the gophkeeper service.
// The command accepts two flags: email and password.
// If the email or password flags are not provided, the command will fail.
//
// The command performs the following steps:
// 1. Retrieves the email and password from the command-line flags.
// 2. Creates a new AuthClient using the provided configuration.
// 3. Calls the Login method of the AuthClient with the email and password.
// 4. If the login is successful, it saves the access token_cache to the token_cache storage.
// 5. Prints the access token_cache to the console.
//
// Parameters:
// - None.
//
// Return:
// - A pointer to the Cobra command object.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Logins a user in the gophkeeper service",
		Run:   loginCmdRun,
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

// loginCmdRun is the main function for the login command. It handles user login to the gophkeeper service.
//
// The function performs the following steps:
// 1. Retrieves the email and password from the command-line flags.
// 2. Creates a new AuthClient using the provided configuration.
// 3. Calls the Login method of the AuthClient with the email and password.
// 4. If the login is successful, it saves the access token_cache to the token_cache storage.
// 5. Prints the access token_cache to the console.
//
// Parameters:
// - cmd: A pointer to the Cobra command object.
// - args: An array of strings representing the command-line arguments.
//
// Return:
// - None.
func loginCmdRun(cmd *cobra.Command, args []string) {
	const op = "client.auth.login.run"
	log := logger.GetInstance().Log.With("op", op)

	email, err := cmd.Flags().GetString("email")
	if err != nil {
		log.Error("Error while getting email", slog.String("error", err.Error()))
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		log.Error("Error while getting password", slog.String("error", err.Error()))
	}

	cfg := config.GetInstance().Config
	authClient := app.NewAuthClient(app.GetAuthConnection(log, cfg.Client))

	token, err := authClient.Login(context.Background(), email, password)
	if err != nil {
		log.Error("Error while login", slog.String("error", err.Error()))
	}

	if err := token_cache.GetInstance().Save(token); err != nil {
		log.Error("Failed to store access token_cache", slog.String("error", err.Error()))
	}

	fmt.Printf("Access Token: %s\n", token)
}
