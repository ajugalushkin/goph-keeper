package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd является командой авторизвции, используется совместно с
// register и login
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage user registration, authentication and authorization",
}

// init метод инициализации, добавляет authCmd к rootCmd
func init() {
	rootCmd.AddCommand(authCmd)
}
