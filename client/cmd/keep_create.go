/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// keepCreateCmd represents the create command
var keepCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
}

func init() {
	keepCmd.AddCommand(keepCreateCmd)
}
