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

// NewAuthClient returns a new auth client
func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	service := authv1.NewAuthServiceV1Client(cc)
	return &AuthClient{service}
}

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
