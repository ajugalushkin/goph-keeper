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
		keeperClient := app.NewKeeperClient(app.GetKeeperConnection(token))

		resp, err := keeperClient.(context.Background(), &v1.CreateItemRequestV1{
			Name:    name,
			Content: content,
		})

		//resp, err := secretClient.UpdateSecret(context.Background(), &pb.UpdateSecretRequest{
		//	Name:    name,
		//	Content: content,
		//})
		//if err != nil {
		//	log.Fatal().Msgf("Failed to update secret: %v", err)
		//	return
		//}

		fmt.Printf("Secret %s version %v updated successfully\n", resp.GetName(), resp.GetVersion())
	},
}

func init() {
	keepUpdateCmd.AddCommand(keepUpdateCardCmd)
}
