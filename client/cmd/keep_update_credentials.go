package cmd

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// credentialsCmd represents the credentials command
var keepUpdateCredentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Update credentials secret",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep update credentials"
		log := logger.GetInstance().Log.With("op", op)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error("Error reading secret name: ",
				slog.String("error", err.Error()))
		}

		login, err := cmd.Flags().GetString("login")
		if err != nil {
			log.Error("Error reading login: ",
				slog.String("error", err.Error()))
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Error("Error reading password: ",
				slog.String("error", err.Error()))
		}

		credentials := vaulttypes.Credentials{
			Login:    login,
			Password: password,
		}

		content, err := encryptSecret(credentials)
		if err != nil {
			log.Error("Failed to encrypt secret: ",
				slog.String("error", err.Error()))
			return
		}

		token, err := tokenStorage.Load()
		if err != nil {
			return
		}
		cfg := config.GetInstance().Config.Client
		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

		resp, err := keeperClient.UpdateItem(context.Background(), &v1.UpdateItemRequestV1{
			Name:    name,
			Content: content,
		})
		if err != nil {
			log.Error("Failed to update secret: ",
				slog.String("error", err.Error()))
			return
		}

		fmt.Printf("Secret %s version %v updated successfully\n", resp.GetName(), resp.GetVersion())
	},
}

func init() {
	const op = "keep update credentials"
	keepUpdateCmd.AddCommand(keepUpdateCredentialsCmd)

	keepUpdateCredentialsCmd.Flags().String("name", "", "Secret name")
	if err := keepUpdateCredentialsCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepUpdateCredentialsCmd.Flags().String("login", "", "Login")
	if err := keepUpdateCredentialsCmd.MarkFlagRequired("login"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepUpdateCredentialsCmd.Flags().String("password", "", "Password")
	if err := keepUpdateCredentialsCmd.MarkFlagRequired("password"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
}
