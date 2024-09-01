/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// keepCmd represents the keep command
var keepCmd = &cobra.Command{
	Use:   "keep",
	Short: "A brief description of your command",
}

func init() {
	rootCmd.AddCommand(keepCmd)
}
