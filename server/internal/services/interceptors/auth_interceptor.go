package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ajugalushkin/goph-keeper/server/internal/lib/jwt"
)

type AuthInterceptor struct {
	jwtManager        *jwt.JWTManager
	accessibleMethods []string
}

func NewAuthInterceptor(jwtManager *jwt.JWTManager, accessibleMethods []string) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, accessibleMethods}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
	isMethodExistsFnc := func(method string) bool {
		for i := 0; i < len(interceptor.accessibleMethods); i++ {
			if method == interceptor.accessibleMethods[i] {
				return true
			}
		}
		return false
	}

	if ok := isMethodExistsFnc(method); !ok {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	ok, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	if ok {
		return nil
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
