package v1

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/auth/v1/mocks"
	"github.com/ajugalushkin/goph-keeper/server/internal/services"
)

// Successfully registers the AuthServiceV1Server with the provided gRPC server
func TestRegister_Success(t *testing.T) {
	// Arrange
	gRPCServer := grpc.NewServer()
	mockAuth := mocks.NewAuth(t)

	// Act
	Register(gRPCServer, mockAuth)

	// Assert
	serviceDesc, ok := gRPCServer.GetServiceInfo()["auth.v1.AuthServiceV1"]
	if !ok {
		t.Fatalf("Expected service to be registered, but it was not")
	}
	if reflect.ValueOf(serviceDesc).IsZero() {
		t.Fatalf("Service description should not be nil")
	}
}

// The gRPC server pointer is nil
func TestRegister_NilGRPCServer(t *testing.T) {
	// Arrange
	var gRPCServer *grpc.Server = nil
	mockAuth := mocks.NewAuth(t)

	// Act and Assert
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic when gRPC server is nil, but did not get one")
		}
	}()

	Register(gRPCServer, mockAuth)
}

// Successfully registers a new user with valid email and password
func TestRegisterV1_Success(t *testing.T) {
	ctx := context.Background()
	req := &v1.RegisterRequestV1{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockAuth := mocks.NewAuth(t)
	mockAuth.On("RegisterNewUser", ctx, req.GetEmail(), req.GetPassword()).Return(int64(1), nil)

	s := &serverAPI{
		auth: mockAuth,
	}

	resp, err := s.RegisterV1(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.UserId)
}

// Returns an error when the registration request is invalid
func TestRegisterV1_InvalidRequest(t *testing.T) {
	ctx := context.Background()
	req := &v1.RegisterRequestV1{
		Email:    "",
		Password: "",
	}

	mockAuth := mocks.NewAuth(t)

	s := &serverAPI{
		auth: mockAuth,
	}

	_, err := s.RegisterV1(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

// Valid login request returns a token
func TestValidLoginRequestReturnsToken(t *testing.T) {
	ctx := context.Background()
	req := &v1.LoginRequestV1{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockAuth := mocks.NewAuth(t)
	mockAuth.On("Login", ctx, "test@example.com", "password123").Return("valid-token", nil)

	s := &serverAPI{
		auth: mockAuth,
	}

	resp, err := s.LoginV1(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, "valid-token", resp.Token)
}

// Invalid credentials return an invalid argument error
func TestInvalidCredentialsReturnInvalidArgumentError(t *testing.T) {
	ctx := context.Background()
	req := &v1.LoginRequestV1{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockAuth := mocks.NewAuth(t)
	mockAuth.On("Login", ctx, "test@example.com", "wrongpassword").Return("", services.ErrInvalidCredentials)

	s := &serverAPI{
		auth: mockAuth,
	}

	resp, err := s.LoginV1(ctx, req)
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), "invalid credentials")
}
