/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage user registration, authentication and authorization",
	//Run: func(cmd *cobra.Command, args []string) {
	//	const op = "client.auth"
	//	log := logger.GetInstance().Log.With("op", op)
	//
	//	if AuthClientConnection == nil {
	//		log.Error("Unable to connect to server")
	//	}
	//
	//	AuthClient = app.NewAuthClient(AuthClientConnection)
	//},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
