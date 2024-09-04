package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// getCmd represents the get command
var keepGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secret",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep_get"
		log := logger.GetInstance().Log.With("op", op)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error("Error reading secret name: ", slog.String("error", err.Error()))
		}

		//resp, err := secretClient.GetSecret(context.Background(), &pb.GetSecretRequest{
		//	Name: name,
		//})
		//if err != nil {
		//	log.Fatal().Err(err).Msg("Failed to get secret")
		//}
		//
		//secret, err := decryptSecret(resp.GetContent())
		//if err != nil {
		//	log.Fatal().Err(err).Msg("Failed to decrypt secret")
		//}

		//fmt.Printf("%s\n", secret)
	},
}

func init() {
	const op = "keep_get"

	keepCmd.AddCommand(keepGetCmd)

	keepGetCmd.Flags().String("name", "", "Secret name")
	if err := keepGetCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
}
