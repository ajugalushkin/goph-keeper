package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// структура для root уровня.
var (
	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "Application GophKeeper",
		Long:  `Application GophKeeper`,
	}
)

// Execute позволяет вызывать root.Execute из другого пакета.
func Execute() error {
	return rootCmd.Execute()
}

// init функция позволяет считать параметры запуска из флагов,
// для чтения используется cobra + viper.
func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringP("ServerAddress", "a", "", "address and port to run server")
	rootCmd.PersistentFlags().StringP("TokenTTL", "t", "", "Token TTL")
	rootCmd.PersistentFlags().StringP("Env", "e", "", "Env dev/prod")

	_ = viper.BindPFlag("Server_Address", rootCmd.PersistentFlags().Lookup("ServerAddress"))
	_ = viper.BindPFlag("Token_TTL", rootCmd.PersistentFlags().Lookup("TokenTTL"))
	_ = viper.BindPFlag("Env", rootCmd.PersistentFlags().Lookup("Env"))

}
