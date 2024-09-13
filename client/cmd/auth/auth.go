package auth

import (
	"github.com/ajugalushkin/goph-keeper/client/cmd/auth/login"
	"github.com/ajugalushkin/goph-keeper/client/cmd/auth/register"
	"github.com/spf13/cobra"
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
func NewCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "auth",
        Short: "Manage user registration, authentication and authorization",
    }
    
    cmd.AddCommand(login.NewCommand())
    cmd.AddCommand(register.NewCommand())

    return cmd
}
