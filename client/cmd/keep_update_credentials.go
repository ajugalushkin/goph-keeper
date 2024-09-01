/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// credentialsCmd represents the credentials command
var keepUpdateCredentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("credentials called")
	},
}

func init() {
	keepUpdateCmd.AddCommand(keepUpdateCredentialsCmd)
}
