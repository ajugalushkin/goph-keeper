package v1

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
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
) (*authv1.RegisterResponseV1, error) {
	v, err := protovalidate.New()
	if err != nil {
		return nil, err
	}
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	return &authv1.RegisterResponseV1{}, nil
}

func LoginV1(
	ctx context.Context,
	req *authv1.LoginRequestV1,
) (*authv1.LoginRequestV1, error) {
	panic("implement me")
}
