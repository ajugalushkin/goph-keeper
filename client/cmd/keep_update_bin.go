/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// binCmd represents the bin command
var keepUpdateBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bin called")
	},
}

func init() {
	keepUpdateCmd.AddCommand(keepUpdateBinCmd)
}
