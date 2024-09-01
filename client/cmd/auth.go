/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	"github.com/ajugalushkin/goph-keeper/client/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

var AuthClient *app.AuthClient

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage user registration, authentication and authorization",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		const op = "client.auth"
		log := logger.GetInstance().Log.With("op", op)

		var err error
		cfg := config.GetInstance().Config

		retryOpts := []grpcretry.CallOption{
			grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
			grpcretry.WithMax(uint(cfg.Client.Retries)),
			grpcretry.WithPerRetryTimeout(cfg.Client.Timeout),
		}

		logOpts := []grpclog.Option{
			grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
		}

		cc, err := grpc.DialContext(context.Background(), cfg.Client.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(
				grpclog.UnaryClientInterceptor(app.InterceptorLogger(log), logOpts...),
				grpcretry.UnaryClientInterceptor(retryOpts...),
			),
		)

		if err != nil {
			log.Error("Unable to connect to server", "error", err)
		}

		AuthClient = app.NewAuthClient(cc)

		interceptor, err := app.NewAuthInterceptor(AuthClient, authMethods(), refreshDuration)
		if err != nil {
			//log.Fatal("cannot create auth interceptor: ", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
