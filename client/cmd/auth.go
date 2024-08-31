/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage user registration, authentication and authorization",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("auth called")
		//c, err := client.New(context.Background(), log, "", "", ``)
		//if err != nil {
		//	return
		//}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
