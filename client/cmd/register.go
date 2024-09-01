/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"log/slog"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers a user in the gophkeeper service",
	//Run: func(cmd *cobra.Command, args []string) {
	//	const op = "cmd.authCmd.auth.register"
	//	env := viper.Get("env")
	//
	//	//log := logger.GetInstance().Log.With("op", op)
	//
	//	log := slog.New(
	//		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	//	)
	//
	//	email, err := cmd.Flags().GetString("email")
	//	if err != nil {
	//		log.Error("Error getting email", "error", err)
	//	}
	//
	//	password, err := cmd.Flags().GetString("password")
	//	if err != nil {
	//		log.Error("Error getting password", "error", err)
	//	}
	//
	//	address, err := cmd.Flags().GetString("address")
	//	if err != nil {
	//		log.Error("Error getting address", "error", err)
	//	}
	//
	//	timeout, err := cmd.Flags().GetDuration("timeout")
	//	if err != nil {
	//		log.Error("Error getting timeout", "error", err)
	//	}
	//
	//	retriesCount, err := cmd.Flags().GetInt("retries_count")
	//	if err != nil {
	//		log.Error("Error getting timeout", "error", err)
	//	}
	//
	//	client, err := app.New(context.Background(), log, address, timeout, retriesCount)
	//	if err != nil {
	//		log.Error("Error creating client", "error", err)
	//	}
	//
	//	err = client.Register(context.Background(), email, password)
	//	if err != nil {
	//		log.Error("Error registering user", "error", err)
	//	}
	//},
}

func init() {
	authCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP("email", "e", "", "User Email")
	if err := registerCmd.MarkFlagRequired("email"); err != nil {
		slog.Error("Error setting email flag", "error", err)
	}
	registerCmd.Flags().StringP("password", "p", "", "User password")
	if err := registerCmd.MarkFlagRequired("password"); err != nil {
		slog.Error("Error setting password flag", "error", err)
	}
}
