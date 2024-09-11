package cmd

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
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
			log.Error("Error getting email", slog.String("error", err.Error()))
			return
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Error("Error getting password", slog.String("error", err.Error()))
			return
		}

		cfg := config.GetInstance().Config
		authClient := app.NewAuthClient(app.GetAuthConnection(log, cfg.Client))

		err = authClient.Register(context.Background(), email, password)
		if err != nil {
			log.Error("Error registering user", slog.String("error", err.Error()))
			return
		}

		fmt.Println("User registered successfully")
	},
}

func init() {
	authCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP("email", "e", "", "User Email")
	if err := registerCmd.MarkFlagRequired("email"); err != nil {
		slog.Error("Error setting email flag", slog.String("error", err.Error()))
	}
	registerCmd.Flags().StringP("password", "p", "", "User password")
	if err := registerCmd.MarkFlagRequired("password"); err != nil {
		slog.Error("Error setting password flag", slog.String("error", err.Error()))
	}
}
