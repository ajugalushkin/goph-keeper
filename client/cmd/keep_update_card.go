/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cardCmd represents the card command
var keepUpdateCardCmd = &cobra.Command{
	Use:   "card",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("card called")
	},
}

func init() {
	keepUpdateCmd.AddCommand(keepUpdateCardCmd)
}
