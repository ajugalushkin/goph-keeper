/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
	Run: func(cmd *cobra.Command, args []string) {
		const op = "client.auth.login.run"
		log := logger.GetInstance().Log.With("op", op)

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Error("Error while getting email", "error", err)
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Error("Error while getting password")
		}

		_, err = AuthClient.Login(context.Background(), email, password)
		if err != nil {
			log.Error("Error while login", "error", err)
		}

		log.Info("Login success")
	},
}

func init() {
	const op = "client.auth.login.init"
	log := logger.GetInstance().Log.With("op", op)

	authCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("email", "e", "", "User Email")
	if err := loginCmd.MarkFlagRequired("email"); err != nil {
		log.Error("Error marking email as required")
	}
	loginCmd.Flags().StringP("password", "p", "", "User password")
	if err := loginCmd.MarkFlagRequired("password"); err != nil {
		log.Error("Error marking password as required")
	}
}
