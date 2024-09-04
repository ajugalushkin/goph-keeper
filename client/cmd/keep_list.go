package cmd

import (
	"github.com/spf13/cobra"
)

// keepListCmd represents the list command
var keepListCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets",
	Run: func(cmd *cobra.Command, args []string) {
		//resp, err := secretClient.ListSecrets(context.Background(), &pb.ListSecretsRequest{})
		//if err != nil {
		//	log.Fatal().Err(err).Msg("Failed to list secret")
		//}
		//
		//for _, info := range resp.GetSecrets() {
		//	secret, err := decryptSecret(info.GetContent())
		//	if err != nil {
		//		log.Fatal().Err(err).Msg("Failed to decrypt secret")
		//	}
		//
		//	fmt.Printf("%s\n", secret)
		//}
	},
}

func init() {
	keepCmd.AddCommand(keepListCmd)
}
