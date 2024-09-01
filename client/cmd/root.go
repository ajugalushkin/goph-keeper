/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"

	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ajugalushkin/goph-keeper/client/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

var cfgFile string

var AuthClientConnection *grpc.ClientConn

var ValtClientConnection *grpc.ClientConn

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gophkeeper_client",
	Short: "GophKeeper cli client",
	Long:  "GophKeeper cli client allows keep and return secrets in/from Keeper server.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		const op = "rootCmd.PersistentPreRun"
		log := logger.GetInstance().Log.With("op", op)

		cfg := config.GetInstance().Config
		retryOpts := []grpcretry.CallOption{
			grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
			grpcretry.WithMax(uint(cfg.Client.Retries)),
			grpcretry.WithPerRetryTimeout(cfg.Client.Timeout),
		}

		logOpts := []grpclog.Option{
			grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
		}

		var err error
		AuthClientConnection, err = grpc.DialContext(context.Background(), cfg.Client.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(
				grpclog.UnaryClientInterceptor(app.InterceptorLogger(log), logOpts...),
				grpcretry.UnaryClientInterceptor(retryOpts...),
			),
		)
		if err != nil {
			log.Error("Unable to connect to server", "error", err)
		}

		interceptor, err := app.NewAuthInterceptor(AuthClient, authMethods())
		if err != nil {
			log.Error("Unable to create interceptor", "error", err)
		}

		ValtClientConnection, err = grpc.DialContext(
			context.Background(),
			cfg.Client.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(interceptor.Unary()),
			grpc.WithStreamInterceptor(interceptor.Stream()),
		)
		if err != nil {
			log.Error("Unable to connect to server", "error", err)
		}

	},
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
		viper.AddConfigPath("./client/config")
		viper.AddConfigPath("/etc/goph-keeper/")
		viper.AddConfigPath("$HOME/.goph-keeper")
		viper.AddConfigPath(".")
	}

	usedFile := viper.ConfigFileUsed()
	slog.Debug("load config file:", usedFile)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			slog.Error("Error reading config file: ", err)
		}
		slog.Debug("Config file not found in ", cfgFile)
	} else {
		slog.Debug("Using config file:", viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		key := strings.ReplaceAll(flag.Name, "-", ".")
		if err := viper.BindPFlag(key, flag); err != nil {
			slog.Error("Error parsing flag: ", err)
		}
	})

	config.GetInstance()
	logger.GetInstance()
}

func authMethods() map[string]bool {
	return map[string]bool{
		keeperv1.KeeperServiceV1_ListItemV1_FullMethodName: true,
		keeperv1.KeeperServiceV1_SetItemV1_FullMethodName:  true,
	}
}
