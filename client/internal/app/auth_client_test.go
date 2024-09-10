package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	authv1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
	"github.com/ajugalushkin/goph-keeper/gen/auth/v1/mocks"
)

// Creates an AuthClient instance with a valid grpc.ClientConn
func TestNewAuthClientWithValidConn(t *testing.T) {
	// Arrange
	conn := new(grpc.ClientConn)

	// Act
	client := NewAuthClient(conn)

	// Assert
	if client == nil {
		t.Errorf("Expected AuthClient instance, got nil")
	}
	if client.api == nil {
		t.Errorf("Expected AuthServiceV1Client instance, got nil")
	}
}

// Handles nil grpc.ClientConn input gracefully
func TestNewAuthClientWithNilConn(t *testing.T) {
	// Arrange
	var conn *grpc.ClientConn = nil

	// Act
	client := NewAuthClient(conn)

	// Assert
	if client == nil {
		t.Errorf("Expected AuthClient instance, got nil")
	}
	if client.api == nil {
		t.Errorf("Expected AuthServiceV1Client instance, got nil")
	}
}

// Successful registration with valid email and password
func TestSuccessfulRegistration(t *testing.T) {
	ctx := context.Background()
	mockAuthService := mocks.NewAuthServiceV1Client(t)
	authClient := &AuthClient{api: mockAuthService}

	mockAuthService.On("RegisterV1", ctx, &authv1.RegisterRequestV1{
		Email:    "test@example.com",
		Password: "password123",
	}).Return(&authv1.RegisterResponseV1{}, nil)

	err := authClient.Register(ctx, "test@example.com", "password123")
	assert.NoError(t, err)
	mockAuthService.AssertExpectations(t)
}

// Registration with empty email
func TestRegistrationWithEmptyEmail(t *testing.T) {
	ctx := context.Background()
	mockAuthService := mocks.NewAuthServiceV1Client(t)
	authClient := &AuthClient{api: mockAuthService}

	mockAuthService.On("RegisterV1", ctx, &authv1.RegisterRequestV1{
		Email:    "",
		Password: "password123",
	}).Return(nil, fmt.Errorf("email cannot be empty"))

	err := authClient.Register(ctx, "", "password123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email cannot be empty")
	mockAuthService.AssertExpectations(t)
}

// Successful login returns a valid token
func TestSuccessfulLoginReturnsValidToken(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	expectedToken := "valid-token"

	mockAuthService := mocks.NewAuthServiceV1Client(t)
	mockAuthService.On("LoginV1", ctx, &authv1.LoginRequestV1{
		Email:    email,
		Password: password,
	}).Return(&authv1.LoginResponseV1{Token: expectedToken}, nil)

	client := &AuthClient{api: mockAuthService}

	token, err := client.Login(ctx, email, password)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

// Login with incorrect email or password returns an error
func TestLoginWithIncorrectCredentialsReturnsError(t *testing.T) {
	ctx := context.Background()
	email := "wrong@example.com"
	password := "wrongpassword"
	expectedError := fmt.Errorf("client.auth.Login: invalid credentials")

	mockAuthService := mocks.NewAuthServiceV1Client(t)
	mockAuthService.On("LoginV1", ctx, &authv1.LoginRequestV1{
		Email:    email,
		Password: password,
	}).Return(nil, expectedError)

	client := &AuthClient{api: mockAuthService}

	token, err := client.Login(ctx, email, password)
	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.Equal(t, expectedError.Error(), err.Error())
}
