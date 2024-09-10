package services

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
)

// Creates an AuthInterceptor instance with valid logger, JWT manager, and accessible methods
func TestNewAuthInterceptor_ValidInputs(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	jwtManager := &JWTManager{
		log:           log,
		secretKey:     "test_secret",
		tokenDuration: time.Hour,
	}
	accessibleMethods := []string{"/test.Method"}

	interceptor := NewAuthInterceptor(log, jwtManager, accessibleMethods)

	if interceptor.log != log {
		t.Errorf("expected logger to be %v, got %v", log, interceptor.log)
	}
	if interceptor.jwtManager != jwtManager {
		t.Errorf("expected JWT manager to be %v, got %v", jwtManager, interceptor.jwtManager)
	}
	if !reflect.DeepEqual(interceptor.accessibleMethods, accessibleMethods) {
		t.Errorf("expected accessible methods to be %v, got %v", accessibleMethods, interceptor.accessibleMethods)
	}
}

// Logs the start of the unary interceptor with method name
func TestLogsStartOfUnaryInterceptorWithMethodName(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	jwtManager := &JWTManager{}
	accessibleMethods := []string{"/test.Service/Method"}
	interceptor := NewAuthInterceptor(log, jwtManager, accessibleMethods)

	ctx := context.Background()
	req := struct{}{}
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	unaryInterceptor := interceptor.Unary()
	_, err := unaryInterceptor(ctx, req, info, handler)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
