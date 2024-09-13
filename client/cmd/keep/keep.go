package keep

import (
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/create"
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/del"
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/get"
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/list"
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep/update"
	"github.com/spf13/cobra"
)


// NewCommand creates a new cobra.Command for managing user private data.
// The command includes subcommands for creating, getting, deleting, listing, and updating private data.
//
// The command structure is as follows:
// - Use: "keep"
// - Short: "Manage user private data"
//
// Subcommands:
// - create: Create a new piece of private data.
// - get: Retrieve a specific piece of private data.
// - del: Delete a specific piece of private data.
// - list: List all private data.
// - update: Update a specific piece of private data.
func NewCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "keep",
        Short: "Manage user private data",
    }

    cmd.AddCommand(create.NewCommand())
    cmd.AddCommand(get.NewCommand())
    cmd.AddCommand(del.NewCommand())
    cmd.AddCommand(list.NewCommand())
    cmd.AddCommand(update.NewCommand())

    return cmd
}
