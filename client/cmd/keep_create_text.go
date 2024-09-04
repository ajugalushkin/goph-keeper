package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// textCmd represents the text command
var keepCreateTextCmd = &cobra.Command{
	Use:   "text",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("text called")
	},
}

func init() {
	keepCreateCmd.AddCommand(keepCreateTextCmd)
}
