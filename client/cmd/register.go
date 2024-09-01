/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers a user in the gophkeeper service",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "client.auth.register.run"
		log := logger.GetInstance().Log.With("op", op)

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Error("Error getting email", "error", err)
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Error("Error getting password", "error", err)
		}

		err = AuthClient.Register(context.Background(), email, password)
		if err != nil {
			log.Error("Error registering user", "error", err)
		}
	},
}

func init() {
	const op = "client.auth.register.init"
	log := logger.GetInstance().Log.With("op", op)

	authCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP("email", "e", "", "User Email")
	if err := registerCmd.MarkFlagRequired("email"); err != nil {
		log.Error("Error setting email flag", "error", err)
	}
	registerCmd.Flags().StringP("password", "p", "", "User password")
	if err := registerCmd.MarkFlagRequired("password"); err != nil {
		log.Error("Error setting password flag", "error", err)
	}
}
