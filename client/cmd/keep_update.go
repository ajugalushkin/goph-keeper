package cmd

import (
	"github.com/spf13/cobra"
)

// keepUpdateCmd represents the update command
var keepUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
}

func init() {
	keepCmd.AddCommand(keepUpdateCmd)
}
