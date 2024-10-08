package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ajugalushkin/goph-keeper/gen/auth/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/services"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

type serverAPI struct {
	v1.UnimplementedAuthServiceV1Server
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	v1.RegisterAuthServiceV1Server(gRPC, &serverAPI{
		auth: auth,
	})
}

func (s *serverAPI) RegisterV1(
	ctx context.Context,
	req *v1.RegisterRequestV1,
) (*v1.RegisterResponseV1, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register new user")
	}

	return &v1.RegisterResponseV1{UserId: user}, nil
}

func (s *serverAPI) LoginV1(
	ctx context.Context,
	req *v1.LoginRequestV1,
) (*v1.LoginResponseV1, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")

		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &v1.LoginResponseV1{
		Token: token,
	}, nil
}
