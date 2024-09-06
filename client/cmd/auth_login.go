package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logins a user in the gophkeeper service",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "client.auth.login.run"
		log := logger.GetInstance().Log.With("op", op)

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Error("Error while getting email", slog.String("error", err.Error()))
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Error("Error while getting password", slog.String("error", err.Error()))
		}

		authClient := app.NewAuthClient(app.GetAuthConnection())

		token, err := authClient.Login(context.Background(), email, password)
		if err != nil {
			log.Error("Error while login", slog.String("error", err.Error()))
		}

		if err := tokenStorage.Save(token); err != nil {
			log.Error("Failed to store access token", slog.String("error", err.Error()))
		}

		fmt.Printf("Access Token: %s\n", token)
	},
}

func init() {
	authCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("email", "e", "", "User Email")
	if err := loginCmd.MarkFlagRequired("email"); err != nil {
		slog.Error("Error marking email as required", slog.String("error", err.Error()))
	}
	loginCmd.Flags().StringP("password", "p", "", "User password")
	if err := loginCmd.MarkFlagRequired("password"); err != nil {
		slog.Error("Error marking password as required", slog.String("error", err.Error()))
	}
}
