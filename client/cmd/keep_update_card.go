package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// keepUpdateCardCmd represents the card command
var keepUpdateCardCmd = &cobra.Command{
	Use:   "card",
	Short: "Update card secret",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep update card"
		log := logger.GetInstance().Log.With("op", op)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error("Error reading secret name: ",
				slog.String("error", err.Error()))
			return
		}

		number, err := cmd.Flags().GetString("number")
		if err != nil {
			log.Error("Error reading card number: ",
				slog.String("error", err.Error()))
			return
		}

		date, err := cmd.Flags().GetString("date")
		if err != nil {
			log.Error("Error reading card expiry date: ",
				slog.String("error", err.Error()))
			return
		}

		code, err := cmd.Flags().GetString("code")
		if err != nil {
			log.Error("Error reading card security code: ",
				slog.String("error", err.Error()))
			return
		}

		holder, err := cmd.Flags().GetString("holder")
		if err != nil {
			log.Error("Error reading card holder: ",
				slog.String("error", err.Error()))
			return
		}

		card := vaulttypes.Card{
			Number:       number,
			ExpiryDate:   date,
			SecurityCode: code,
			Holder:       holder,
		}

		content, err := encryptSecret(card)
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
	const op = "keep update card"
	keepUpdateCmd.AddCommand(keepUpdateCardCmd)

	keepUpdateCardCmd.Flags().String("name", "", "Secret name")
	if err := keepUpdateCardCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepUpdateCardCmd.Flags().String("number", "", "Card number")
	if err := keepUpdateCardCmd.MarkFlagRequired("number"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepUpdateCardCmd.Flags().String("date", "", "Card expiry date")
	if err := keepUpdateCardCmd.MarkFlagRequired("date"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepUpdateCardCmd.Flags().String("code", "", "Card security code")
	if err := keepUpdateCardCmd.MarkFlagRequired("code"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepUpdateCardCmd.Flags().String("holder", "", "Card holder")
	if err := keepUpdateCardCmd.MarkFlagRequired("holder"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
}
