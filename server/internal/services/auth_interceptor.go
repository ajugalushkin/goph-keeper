package services

import (
	"context"
	"log"
	"log/slog"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type key int

const (
	// ContextKeyUserID ключ для добавления UserID в контекст при аутентификации
	ContextKeyUserID key = iota
)

type AuthInterceptor struct {
	log               *slog.Logger
	jwtManager        TokenManager
	accessibleMethods []string
}

// NewAuthInterceptor creates a new instance of AuthInterceptor.
// The AuthInterceptor is a gRPC server interceptor that handles authentication and authorization for RPC methods.
//
// Parameters:
// - log: A pointer to a slog.Logger instance for logging.
// - jwtManager: A TokenManager instance for managing JWT tokens.
// - accessibleMethods: A slice of strings representing the RPC methods that are accessible without authentication.
//
// Returns:
// - A pointer to a new AuthInterceptor instance.
func NewAuthInterceptor(
	log *slog.Logger,
	jwtManager TokenManager,
	accessibleMethods []string,
) *AuthInterceptor {
	return &AuthInterceptor{
		log:               log,
		jwtManager:        jwtManager,
		accessibleMethods: accessibleMethods,
	}
}

// Unary returns a unary server interceptor for the AuthInterceptor.
// This interceptor handles authentication and authorization for unary RPC methods.
// It checks if the RPC method is accessible without authentication or if it requires authentication.
// If authentication is required, it verifies the JWT token_cache provided in the metadata.
// If the token_cache is valid, it adds the UserID to the context.
// If the token_cache is invalid or the user does not have permission to access the RPC method, it returns an error.
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		const op = "interceptors.AuthInterceptor.Unary"
		log := interceptor.log.With("op", op)

		log.Info("--> unary interceptor: ", "method", info.FullMethod)

		newCtx, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			log.Debug("unauthorized access method: ", "method", info.FullMethod)
			return nil, err
		}

		return handler(newCtx, req)
	}
}

// Stream returns a stream server interceptor for the AuthInterceptor.
// This interceptor handles authentication and authorization for stream RPC methods.
// It checks if the RPC method is accessible without authentication or if it requires authentication.
// If authentication is required, it verifies the JWT token_cache provided in the metadata.
// If the token_cache is valid, it adds the UserID to the context.
// If the token_cache is invalid or the user does not have permission to access the RPC method, it returns an error.
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> stream interceptor: ", info.FullMethod)

		newCtx, err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, &grpc_middleware.WrappedServerStream{
			ServerStream:   stream,
			WrappedContext: newCtx,
		})
	}
}

// authorize is a helper function for the AuthInterceptor that handles authentication and authorization.
// It checks if the RPC method is accessible without authentication or if it requires authentication.
// If authentication is required, it verifies the JWT token_cache provided in the metadata.
// If the token_cache is valid, it adds the UserID to the context.
// If the token_cache is invalid or the user does not have permission to access the RPC method, it returns an error.
//
// Parameters:
// - ctx: The context for the RPC call.
// - method: The full method name of the RPC call.
//
// Returns:
// - A new context with the UserID added if the token_cache is valid.
// - An error if the token_cache is invalid or the user does not have permission to access the RPC method.
func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
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
		log.Debug("authorized access method: ", "method", method)
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Debug("metadata is empty")
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		log.Debug("token_cache is empty")
		return nil, status.Errorf(codes.Unauthenticated, "authorization token_cache is not provided")
	}

	accessToken := strings.TrimPrefix(values[0], "Bearer ")
	ok, userID, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		log.Debug("invalid access token_cache: ",
			"token_cache", accessToken,
			"values", values)
		return nil, status.Errorf(codes.Unauthenticated, "access token_cache is invalid: %v", err)
	}

	if ok {
		log.Debug("authorized access token_cache: ", slog.String("token_cache", accessToken))
		return context.WithValue(ctx, ContextKeyUserID, userID), nil
	}

	return nil, status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
