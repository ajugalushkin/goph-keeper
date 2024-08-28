package app

import (
	v1 "github.com/ajugalushkin/goph-keeper/common/gen/keeper/v1"
)

type Client struct {
	//	api v1.KeeperServiceV1Client
	api v1.KeeperServiceV1Client
}

//
//func New(c
//	ctx context.Context,
//	addr string,
//	timeout time.Duration,
//	retriesCount int,
//) (*Client, error) {
//	const op = "grpc.New"
//
//	retryOpts := []grpcretry.CallOption{
//		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
//		grpcretry.WithMax(uint(retriesCount)),
//		grpcretry.WithPerRetryTimeout(timeout),
//	}
//
//	//logOpts := []grpclog.Option{
//	//	grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
//	//}
//
//	cc, err := grpc.NewClient(addr,
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//		grpc.WithChainUnaryInterceptor(
//			//grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
//			grpcretry.UnaryClientInterceptor(retryOpts...),
//		),
//	)
//
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	return &Client{
//		//api: v1.NewAuthV1Client(cc),
//	}, nil
//}
//
//func InterceptorLogger(l *slog.Logger) grpclog.Logger {
//	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
//		l.Log(ctx, slog.Level(lvl), msg, fields...)
//	})
//}
//
//func (c *Client) Register(ctx context.Context, email string, password string) error {
//	const op = "client.keeper.Register"
//
//	_, err := c.api.RegisterV1(ctx, &v1.RegisterRequestV1{
//		Email:    email,
//		Password: password,
//	})
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//
//	return nil
//}
//
//func (c *Client) Login(ctx context.Context, email string, password string) (string, error) {
//	const op = "client.keeper.Login"
//
//	resp, err := c.api.LoginV1(ctx, &v1.LoginRequestV1{
//		Email:    email,
//		Password: password,
//	})
//	if err != nil {
//		return "", fmt.Errorf("%s: %w", op, err)
//	}
//
//	return resp.Token, nil
//}
