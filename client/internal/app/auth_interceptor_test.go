package app

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Initializes AuthInterceptor with valid token and authMethods
func TestInitializesAuthInterceptorWithValidTokenAndAuthMethods(t *testing.T) {
	token := "valid-token"
	authMethods := map[string]bool{
		"/service/method": true,
	}

	interceptor, err := NewAuthInterceptor(token, authMethods)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if interceptor.accessToken != token {
		t.Errorf("expected accessToken to be %s, got %s", token, interceptor.accessToken)
	}

	if len(interceptor.authMethods) != 1 || !interceptor.authMethods["/service/method"] {
		t.Errorf("expected authMethods to contain /service/method")
	}
}

// Handles empty token string gracefully
func TestHandlesEmptyTokenStringGracefully(t *testing.T) {
	token := ""
	authMethods := map[string]bool{
		"/service/method": true,
	}

	interceptor, err := NewAuthInterceptor(token, authMethods)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if interceptor.accessToken != token {
		t.Errorf("expected accessToken to be empty, got %s", interceptor.accessToken)
	}

	if len(interceptor.authMethods) != 1 || !interceptor.authMethods["/service/method"] {
		t.Errorf("expected authMethods to contain /service/method")
	}
}

// Interceptor adds authorization metadata to context
func TestInterceptorAddsAuthorizationMetadata(t *testing.T) {
	authMethods := map[string]bool{"TestMethod": true}
	interceptor, err := NewAuthInterceptor("test-token", authMethods)
	if err != nil {
		t.Fatalf("Failed to create AuthInterceptor: %v", err)
	}

	ctx := context.Background()
	method := "TestMethod"
	req, reply := struct{}{}, struct{}{}
	cc := &grpc.ClientConn{}
	invoker := func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			t.Fatalf("No metadata in context")
		}
		if len(md["authorization"]) == 0 || md["authorization"][0] != "Bearer test-token" {
			t.Fatalf("Authorization token not found in metadata")
		}
		return nil
	}

	err = interceptor.Unary()(ctx, method, req, reply, cc, invoker)
	if err != nil {
		t.Fatalf("Unary interceptor returned an error: %v", err)
	}
}

// Method name is empty
func TestInterceptorWithEmptyMethodName(t *testing.T) {
	authMethods := map[string]bool{"": true}
	interceptor, err := NewAuthInterceptor("test-token", authMethods)
	if err != nil {
		t.Fatalf("Failed to create AuthInterceptor: %v", err)
	}

	ctx := context.Background()
	method := ""
	req, reply := struct{}{}, struct{}{}
	cc := &grpc.ClientConn{}
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			t.Fatalf("No metadata in context")
		}
		if len(md["authorization"]) == 0 || md["authorization"][0] != "Bearer test-token" {
			t.Fatalf("Authorization token not found in metadata")
		}
		return nil
	}

	err = interceptor.Unary()(ctx, method, req, reply, cc, invoker)
	if err != nil {
		t.Fatalf("Unary interceptor returned an error: %v", err)
	}
}

// Stream interceptor adds authorization metadata to context
func TestStreamInterceptorAddsAuthorizationMetadata(t *testing.T) {
	authMethods := map[string]bool{"testMethod": true}
	interceptor, err := NewAuthInterceptor("testToken", authMethods)
	if err != nil {
		t.Fatalf("Failed to create AuthInterceptor: %v", err)
	}

	ctx := context.Background()
	desc := &grpc.StreamDesc{}
	cc := &grpc.ClientConn{}
	method := "testMethod"
	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			t.Fatalf("No metadata found in context")
		}
		if len(md["authorization"]) == 0 || md["authorization"][0] != "Bearer testToken" {
			t.Fatalf("Authorization token not found in metadata")
		}
		return nil, nil
	}

	_, err = interceptor.Stream()(ctx, desc, cc, method, streamer)
	if err != nil {
		t.Fatalf("Stream interceptor failed: %v", err)
	}
}
