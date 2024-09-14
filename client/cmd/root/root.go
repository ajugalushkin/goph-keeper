package root

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/ajugalushkin/goph-keeper/client/cmd/auth"
	"github.com/ajugalushkin/goph-keeper/client/cmd/keep"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"

	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

var cfgFile string

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

// NewCommand creates and initializes a new root command for the GophKeeper client.
// The root command is responsible for managing authentication and secret keeping operations.
// It includes subcommands for authentication and secret keeping.
//
// The root command supports a persistent flag for specifying the client configuration file path.
// The flag can be accessed using "--config" or "-c".
//
// The function initializes the token cache storage, reads the client configuration from a file,
// environment variables, and sets up logging.
//
// The function returns a pointer to the initialized root command.
//func NewCommand() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "gophkeeper_client",
//		Short: "GophKeeper cli client",
//		Long:  "GophKeeper cli client allows keep and return secrets in/from Keeper server.",
//	}
//
//	cmd.AddCommand(auth.NewCommand())
//	cmd.AddCommand(keep.NewCommand())
//
//	return cmd
//}

// init initializes the root command and its persistent flags.
// It also sets up the configuration file, environment variables, and logging.
func init() {
	rootCmd.AddCommand(auth.NewCommand())
	rootCmd.AddCommand(keep.NewCommand())

	// Initialize token_cache storage with a file-based storage using "token_cache.txt" as the file path.
	token_cache.GetToken()

	// Add a persistent flag to the root command.
	// The flag is named "config" and can be accessed using "--config" or "-c".
	// It represents the client config filepath.
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "", "Client config filepath")

	// Register the initConfig function to be called before the root command is executed.
	cobra.OnInitialize(initConfig)
}

// initConfig initializes the client configuration by reading from a file, environment variables,
// and setting up logging. It is called before the root command is executed.
func initConfig() {
	if cfgFile == "" {
		cfgFile = os.Getenv("CLIENT_CONFIG")
	}

	// If a config file is specified, use it. Otherwise, search for default config files.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("./client/config")
		viper.AddConfigPath(".")
	}

	// Attempt to read the configuration file.
	if err := viper.ReadInConfig(); err != nil {
		// If the error is not a ConfigFileNotFoundError, log it as an error.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			slog.Error("Error reading config file: ", slog.String("error", err.Error()))
		}
		// Log a message indicating that the config file was not found.
		slog.Info("Config file not found in ", slog.String("file", cfgFile))
	} else {
		// Log a message indicating that the config file was successfully used.
		slog.Info("Using config file: ", slog.String("file", viper.ConfigFileUsed()))
	}

	// Enable automatic population of environment variables.
	viper.AutomaticEnv()
	// Replace hyphens with underscores and periods with underscores in environment variable keys.
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// Bind each flag to its corresponding configuration key.
	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		key := strings.ReplaceAll(flag.Name, "-", ".")
		if err := viper.BindPFlag(key, flag); err != nil {
			slog.Error("Error parsing flag: ", slog.String("error", err.Error()))
		}
	})

	// Initialize the configuration and logger instances.
	config.GetConfig()
	logger.GetLogger()
}
