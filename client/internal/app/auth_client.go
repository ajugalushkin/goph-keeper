package app

import "context"

//go:generate mockery --name AuthClient
type AuthClient interface {
	Register(ctx context.Context, email string, password string) error
	Login(ctx context.Context, email string, password string) (string, error)
}
