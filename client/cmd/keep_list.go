package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// keepListCmd represents the list command
var keepListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	keepCmd.AddCommand(keepListCmd)
}
