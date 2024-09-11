package cmd

import (
	"errors"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/ajugalushkin/goph-keeper/client/internal/token"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

var cfgFile string

var tokenStorage token.Storage

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gophkeeper_client",
	Short: "GophKeeper cli client",
	Long:  "GophKeeper cli client allows keep and return secrets in/from Keeper server.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	tokenStorage = token.NewFileStorage("token.txt")

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "", "Client config filepath")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("./client/config")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			slog.Error("Error reading config file: ", slog.String("error", err.Error()))
		}
		slog.Info("Config file not found in ", slog.String("file", cfgFile))
	} else {
		slog.Info("Using config file: ", slog.String("file", viper.ConfigFileUsed()))
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		key := strings.ReplaceAll(flag.Name, "-", ".")
		if err := viper.BindPFlag(key, flag); err != nil {
			slog.Error("Error parsing flag: ", slog.String("error", err.Error()))
		}
	})

	config.GetInstance()
	logger.GetInstance()
}
