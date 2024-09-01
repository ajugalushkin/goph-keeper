package app

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	authv1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
)

type AuthClient struct {
	api authv1.AuthServiceV1Client
}

//func NewAuthClient(
//	ctx context.Context,
//	log *slog.Logger,
//	addr string,
//	timeout time.Duration,
//	retriesCount int,
//) (*AuthClient, error) {
//	const op = "grpc.New"
//
//	retryOpts := []grpcretry.CallOption{
//		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
//		grpcretry.WithMax(uint(retriesCount)),
//		grpcretry.WithPerRetryTimeout(timeout),
//	}
//
//	logOpts := []grpclog.Option{
//		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
//	}
//
//	cc, err := grpc.DialContext(ctx, addr,
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//		grpc.WithChainUnaryInterceptor(
//			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
//			grpcretry.UnaryClientInterceptor(retryOpts...),
//		),
//	)
//
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	return &AuthClient{
//		api: authv1.NewAuthServiceV1Client(cc),
//	}, nil
//}

//type AuthClient struct {
//	service  pb.AuthServiceClient
//	username string
//	password string
//}

// NewAuthClient returns a new auth client
func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	service := authv1.NewAuthServiceV1Client(cc)
	return &AuthClient{service}
}

//func InterceptorLogger(l *slog.Logger) grpclog.Logger {
//	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
//		l.Log(ctx, slog.Level(lvl), msg, fields...)
//	})
//}

func (c *AuthClient) Register(ctx context.Context, email string, password string) error {
	const op = "client.keeper.Register"

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
	const op = "client.keeper.Login"

	resp, err := c.api.LoginV1(ctx, &authv1.LoginRequestV1{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Token, nil
}
