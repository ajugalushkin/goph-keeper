package cmd

import (
	"github.com/spf13/cobra"
)

// keepCreateCmd represents the create command
var keepCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create secret",
}

func init() {
	keepCmd.AddCommand(keepCreateCmd)
}
