package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gophkeeper_client",
	Short: "GophKeeper cli client",
	Long:  "GophKeeper cli client allows keep and return secrets in/from Keeper server.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to [your-cli-app-name]! Use --help for usage.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
