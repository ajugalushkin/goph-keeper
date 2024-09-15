package update

import (
	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/update/card"
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/update/creds"
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/update/text"
)

// NewCommand creates a new cobra.Command for the "update" subcommand.
// This command is intended to perform an update operation.
//
// The command has the following properties:
// - Use: "update"
// - Short: A brief description of your command
//
// The function returns a pointer to the created cobra.Command.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update secret",
	}

	cmd.AddCommand(card.NewCommand())
	cmd.AddCommand(creds.NewCommand())
	cmd.AddCommand(text.NewCommand())

	return cmd
}
