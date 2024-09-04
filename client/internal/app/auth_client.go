package app

import (
	"context"
	"fmt"
	"log/slog"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ajugalushkin/goph-keeper/client/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	authv1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
)

type AuthClient struct {
	api authv1.AuthServiceV1Client
}

// NewAuthClient returns a new auth client
func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	service := authv1.NewAuthServiceV1Client(cc)
	return &AuthClient{service}
}

func (c *AuthClient) Register(ctx context.Context, email string, password string) error {
	const op = "client.auth.Register"

	_, err := c.api.RegisterV1(ctx, &authv1.RegisterRequestV1{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *AuthClient) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "client.auth.Login"

	resp, err := c.api.LoginV1(ctx, &authv1.LoginRequestV1{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Token, nil
}

func GetAuthConnection() *grpc.ClientConn {
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
	connection, err := grpc.NewClient(cfg.Client.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(interceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		slog.Error("Unable to connect to server", "error", err)
	}

	return connection
}

func interceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
