package v1

import (
	"context"

	"google.golang.org/grpc"

	authv1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
)

type serverAPI struct {
	authv1.UnimplementedAuthV1Server
}

func Register(gRPC *grpc.Server) {
	authv1.RegisterAuthV1Server(gRPC, &serverAPI{})
}

func RegisterV1(
	ctx context.Context,
	req *authv1.RegisterRequestV1,
) (*authv1.RegisterRequestV1, error) {
	panic("implement me")
}

func LoginV1(
	ctx context.Context,
	req *authv1.LoginRequestV1,
) (*authv1.LoginRequestV1, error) {
	panic("implement me")
}
