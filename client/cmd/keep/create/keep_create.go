package create

import (
	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/create/bin"
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/create/card"
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/create/creds"
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/create/text"
)

// NewCommand creates a new cobra.Command for the "create" subcommand.
// This command is responsible for creating different types of secrets, such as binary, card, credentials, and text.
//
// The command structure is as follows:
// - Use: "create"
// - Short: "Create secret"
//
// It adds four subcommands:
// - bin.NewCommand()
// - card.NewCommand()
// - creds.NewCommand()
// - text.NewCommand()
//
// The function returns a pointer to the created cobra.Command.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create secret",
	}

	cmd.AddCommand(bin.NewCommand())
	cmd.AddCommand(card.NewCommand())
	cmd.AddCommand(creds.NewCommand())
	cmd.AddCommand(text.NewCommand())

	return cmd
}
