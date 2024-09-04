package cmd

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// credentialsCmd represents the credentials command
var keepCreateCredentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep create credentials"
		log := logger.GetInstance().Log.With("op", op)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error("Unable to get `name` arg: ", slog.String("error", err.Error()))
			return
		}

		login, err := cmd.Flags().GetString("login")
		if err != nil {
			log.Error("Unable to get `login` arg: ", slog.String("error", err.Error()))
			return
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Error("Unable to get `password` arg: ", slog.String("error", err.Error()))
			return
		}

		credentials := vaulttypes.Credentials{
			Login:    login,
			Password: password,
		}

		content, err := encryptVault(credentials)
		if err != nil {
			log.Error("Failed to encrypt secret: ", slog.String("error", err.Error()))
			return
		}

		token, err := tokenStorage.Load()
		if err != nil {
			return
		}
		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(token))

		_, err = keeperClient.CreateItem(context.Background(), &v1.CreateItemRequestV1{
			Name:    name,
			Content: content,
		})
		if err != nil {
			log.Error("Error while login: ", slog.String("error", err.Error()))
		}

		log.Info("Successfully created credentials")
	},
}

func init() {
	keepCreateCmd.AddCommand(keepCreateCredentialsCmd)

	keepCreateCredentialsCmd.Flags().String("name", "", "Secret name")
	if err := keepCreateCredentialsCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Unable to mark 'name' flag as required %s", slog.String("error", err.Error()))
	}
	keepCreateCredentialsCmd.Flags().String("login", "", "Login")
	if err := keepCreateCredentialsCmd.MarkFlagRequired("login"); err != nil {
		slog.Error("Unable to mark 'login' flag as required %s", slog.String("error", err.Error()))
	}
	keepCreateCredentialsCmd.Flags().String("password", "", "Password")
	if err := keepCreateCredentialsCmd.MarkFlagRequired("password"); err != nil {
		slog.Error("Unable to mark 'password' flag as required %s", slog.String("error", err.Error()))
	}
}
