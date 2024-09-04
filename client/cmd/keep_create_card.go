package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// keepCreateCardCmd represents the card command
var keepCreateCardCmd = &cobra.Command{
	Use:   "card",
	Short: "Create card secret",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep_create_card"
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

		content, err := encryptVault(card)
		if err != nil {
			log.Error("Failed to encrypt secret: ",
				slog.String("error", err.Error()))
			return
		}

		token, err := tokenStorage.Load()
		if err != nil {
			return
		}
		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(token))

		resp, err := keeperClient.CreateItem(context.Background(), &v1.CreateItemRequestV1{
			Name:    name,
			Content: content,
		})
		if err != nil {
			log.Error("Failed to create secret: ", slog.String("error", err.Error()))
		}

		fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())
	},
}

func init() {
	const op = "keep_create_card"
	keepCreateCmd.AddCommand(keepCreateCardCmd)

	keepCreateCardCmd.Flags().String("name", "", "Secret name")
	if err := keepCreateCardCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepCreateCardCmd.Flags().String("number", "", "Card number")
	if err := keepCreateCardCmd.MarkFlagRequired("number"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepCreateCardCmd.Flags().String("date", "", "Card expiry date")
	if err := keepCreateCardCmd.MarkFlagRequired("date"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepCreateCardCmd.Flags().String("code", "", "Card security code")
	if err := keepCreateCardCmd.MarkFlagRequired("code"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	keepCreateCardCmd.Flags().String("holder", "", "Card holder")
	if err := keepCreateCardCmd.MarkFlagRequired("holder"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
}
