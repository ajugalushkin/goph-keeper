package auth

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/cmd/auth/login"
	"github.com/ajugalushkin/goph-keeper/client/cmd/auth/register"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
)

// NewCommand creates a new cobra.Command for managing user authentication and authorization.
//
// The command includes two subcommands:
// - login: Handles user login.
// - register: Handles user registration.
//
// The command is configured with the following properties:
// - Use: "auth"
// - Short: "Manage user registration, authentication and authorization"
//
// The function returns a pointer to the created cobra.Command.
func NewCommand(log *slog.Logger, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage user registration, authentication and authorization",
	}

	cmd.AddCommand(login.NewCommand(log, cfg))
	cmd.AddCommand(register.NewCommand())

	return cmd
}
