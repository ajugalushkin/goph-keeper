package cmd

import (
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/spf13/cobra"
	"log/slog"
)

// deleteCmd represents the delete command
var keepDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		const op = "keep_delete"
		log := logger.GetInstance().Log.With("op", op)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error("Error reading secret name", slog.String("error", err.Error()))
		}

		//resp, err := secretClient.DeleteSecret(
		//	context.Background(), &pb.DeleteSecretRequest{Name: name})
		//if err != nil {
		//	log.Fatal().Msgf("Failed to delete secret: %v", err)
		//	return
		//}
		//
		//fmt.Printf("Secret %s deleted successfully\n", resp.GetName())
	},
}

func init() {
	keepCmd.AddCommand(keepDeleteCmd)
}
