package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ajugalushkin/goph-keeper/server/internal/lib/jwt"
)

type AuthInterceptor struct {
	log               *slog.Logger
	jwtManager        *jwt.JWTManager
	accessibleMethods []string
}

func NewAuthInterceptor(log *slog.Logger, jwtManager *jwt.JWTManager, accessibleMethods []string) *AuthInterceptor {
	return &AuthInterceptor{
		log:               log,
		jwtManager:        jwtManager,
		accessibleMethods: accessibleMethods,
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		const op = "interceptors.AuthInterceptor.Unary"
		log := interceptor.log.With("op", op)

		log.Info("--> unary interceptor: ", info.FullMethod)

		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			log.Debug("unauthorized access method: ", info.FullMethod)
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
	const op = "interceptors.AuthInterceptor.authorize"
	log := interceptor.log.With("op", op)

	isMethodExistsFnc := func(method string) bool {
		for i := 0; i < len(interceptor.accessibleMethods); i++ {
			if method == interceptor.accessibleMethods[i] {
				return true
			}
		}
		return false
	}

	if ok := isMethodExistsFnc(method); ok {
		log.Debug("authorized access method: ", method)
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Debug("metadata is empty")
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		log.Debug("token is empty")
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	ok, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		log.Debug("invalid access token: ", accessToken)
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	if ok {
		log.Debug("authorized access token: ", accessToken)
		return nil
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
