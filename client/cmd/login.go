/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

//// loginCmd represents the login command
//var loginCmd = &cobra.Command{
//	Use:   "login",
//	Short: "A brief description of your command",
//	Long:  `A longer description that spans multiple lines and likely contains examples`,
//	Run: func(cmd *cobra.Command, args []string) {
//		const op = "cmd.authCmd.auth.login"
//		log := logger.GetInstance().Log.With("op", op)
//
//		email, err := cmd.Flags().GetString("email")
//		if err != nil {
//			//log.Fatal().Err(err).Msg("Failed to read email")
//			log.Debug("")
//		}
//
//		password, err := cmd.Flags().GetString("password")
//		if err != nil {
//			//log.Fatal().Err(err).Msg("Failed to read password")
//			log.Debug("")
//		}
//
//		cfg := config.GetInstance().Config
//
//		client, err := app.New(context.Background(), log, cfg.Client.Address, cfg.Client.Timeout, cfg.Client.RetriesCount)
//		if err != nil {
//			//log.Debug("")
//		}
//
//		_, err = client.Login(context.Background(), email, password)
//		if err != nil {
//			//log.Debug("error while registering user", "error", err)
//		}
//	},
//}
//
//func init() {
//	authCmd.AddCommand(registerCmd)
//
//	registerCmd.Flags().StringP("email", "e", "", "User Email")
//	if err := registerCmd.MarkFlagRequired("email"); err != nil {
//		//log.Error().Err(err)
//	}
//	registerCmd.Flags().StringP("password", "p", "", "User password")
//	if err := registerCmd.MarkFlagRequired("password"); err != nil {
//		//log.Error().Err(err)
//	}
//}
