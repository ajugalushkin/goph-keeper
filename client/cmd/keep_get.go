package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var keepGetCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
	},
}

func init() {
	keepCmd.AddCommand(keepGetCmd)
}
