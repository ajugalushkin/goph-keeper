package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// textCmd represents the text command
var keepUpdateTextCmd = &cobra.Command{
	Use:   "text",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("text called")
	},
}

func init() {
	keepUpdateCmd.AddCommand(keepUpdateTextCmd)
}
