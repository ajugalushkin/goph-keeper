package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage user registration, authentication and authorization",
}

func init() {
	rootCmd.AddCommand(authCmd)
}
