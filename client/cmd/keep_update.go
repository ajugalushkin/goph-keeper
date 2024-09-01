/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// keepUpdateCmd represents the update command
var keepUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
	},
}

func init() {
	keepCmd.AddCommand(keepUpdateCmd)
}
