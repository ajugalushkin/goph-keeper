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

//go:generate mockery --name Auth
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

// Register registers the gRPC server for the AuthServiceV1 with the provided gRPC server and Auth implementation.
//
// gRPC: The gRPC server to register the AuthServiceV1 server.
// auth: An implementation of the Auth interface that provides the necessary functionality for the AuthServiceV1.
func Register(gRPC *grpc.Server, auth Auth) {
    v1.RegisterAuthServiceV1Server(gRPC, &serverAPI{
        auth: auth,
    })
}

// RegisterV1 handles the registration of a new user.
//
// It accepts a gRPC context and a RegisterRequestV1 containing the user's email and password.
// The function performs the following steps:
// 1. Validates the input using protovalidate.
// 2. Calls the RegisterNewUser method of the Auth interface to register the new user.
// 3. Returns a RegisterResponseV1 containing the user's ID if successful, or an error if any step fails.
//
// Parameters:
// ctx (context.Context): The gRPC context for the request.
// req (*v1.RegisterRequestV1): The request containing the user's email and password.
//
// Return:
// (*v1.RegisterResponseV1, error): A pointer to the RegisterResponseV1 containing the user's ID, or an error if any step fails.
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

// LoginV1 handles the login process for an existing user.
//
// It accepts a gRPC context and a LoginRequestV1 containing the user's email and password.
// The function performs the following steps:
// 1. Validates the input using protovalidate.
// 2. Calls the Login method of the Auth interface to authenticate the user.
// 3. Returns a LoginResponseV1 containing the user's token if successful, or an error if any step fails.
//
// Parameters:
// ctx (context.Context): The gRPC context for the request.
// req (*v1.LoginRequestV1): The request containing the user's email and password.
//
// Return:
// (*v1.LoginResponseV1, error): A pointer to the LoginResponseV1 containing the user's token, or an error if any step fails.
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
