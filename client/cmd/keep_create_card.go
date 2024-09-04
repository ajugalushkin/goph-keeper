package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cardCmd represents the card command
var keepCreateCardCmd = &cobra.Command{
	Use:   "card",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("card called")
	},
}

func init() {
	keepCreateCmd.AddCommand(keepCreateCardCmd)
}
