package auth

import (
	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/cmd/auth/login"
	"github.com/ajugalushkin/goph-keeper/client/cmd/auth/register"
)

var auth = &cobra.Command{
	Use:   "auth",
	Short: "Manage user registration, authentication and authorization",
}

// NewCommand creates a new cobra.Command for managing user authentication, registration, and authorization.
//
// The function takes two parameters:
// - log: A pointer to a slog.Logger instance for logging.
// - cfg: A pointer to a config.Config instance containing configuration settings.
//
// The function returns a pointer to a cobra.Command with the following properties:
// - Use: Set to "auth".
// - Short: Set to "Manage user registration, authentication and authorization".
//
// The function also initializes an authClient using the app.NewAuthClient function, passing the result of
// app.GetAuthConnection(log, cfg.Client) as an argument.
//
// Finally, the function adds two subcommands to the returned cobra.Command:
// - login.NewCommand(log, authClient)
// - register.NewCommand()
func NewCommand() *cobra.Command {
	return auth
}

func init() {
	auth.AddCommand(login.NewCommand())
	auth.AddCommand(register.NewCommand())
}
