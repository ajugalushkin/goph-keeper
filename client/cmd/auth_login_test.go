package cmd

import (
	"context"
	"errors"
	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

// Command is successfully added to authCmd
func TestCommandAddedToAuthCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "root"}
	authCmd := &cobra.Command{Use: "auth"}
	loginCmd := &cobra.Command{Use: "login"}

	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)

	found := false
	for _, cmd := range authCmd.Commands() {
		if cmd.Use == "login" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("loginCmd was not added to authCmd")
	}
}

// Error occurs when marking email flag as required
func TestErrorMarkingEmailFlagRequired(t *testing.T) {
	loginCmd := &cobra.Command{Use: "login"}
	loginCmd.Flags().StringP("email", "e", "", "User Email")

	err := loginCmd.MarkFlagRequired("email")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Simulate an error scenario by trying to mark a non-existent flag as required
	err = loginCmd.MarkFlagRequired("nonexistent")
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

type MockAuthClient struct {
	mock.Mock
}

func (m *MockAuthClient) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

type MockTokenStorage struct {
	mock.Mock
}

func (m *MockTokenStorage) Save(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func TestLoginWithValidEmailAndPassword(t *testing.T) {
	email := "test@example.com"
	password := "testpassword"

	mockAuthClient := &MockAuthClient{}
	mockAuthClient.On("Login", context.Background(), email, password).Return("test_token", nil)

	token, err := mockAuthClient.Login(context.Background(), email, password)
	assert.Nil(t, err, "Expected no error, but got %v", err)
	assert.Equal(t, "test_token", token, "Expected token 'test_token', but got %s", token)

	mockTokenStorage := &MockTokenStorage{}
	mockTokenStorage.On("Save", token).Return(nil)

	err = mockTokenStorage.Save(token)
	assert.Nil(t, err, "Expected no error, but got %v", err)

	mockAuthClient.AssertExpectations(t)
	mockTokenStorage.AssertExpectations(t)
}

func TestLoginWithInvalidEmailFormat(t *testing.T) {
	const op = "client.auth.login.run"
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	log.With("op", op)

	email := "invalid_email"
	password := "testpassword"

	cfg := config.GetInstance().Config
	authClient := app.NewAuthClient(app.GetAuthConnection(log, cfg.Client))

	_, err := authClient.Login(context.Background(), email, password)
	assert.Error(t, err, "Expected error for invalid email format")
}
func TestLoginWithEmptyEmail(t *testing.T) {
	const op = "client.auth.login.run"
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	log.With("op", op)

	email := ""
	password := "testpassword"

	cfg := config.GetInstance().Config
	authClient := app.NewAuthClient(app.GetAuthConnection(log, cfg.Client))

	_, err := authClient.Login(context.Background(), email, password)
	assert.Error(t, err, "Expected error for empty email")
}
func TestLoginWithEmptyPassword(t *testing.T) {
	const op = "client.auth.login.run"
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	log.With("op", op)

	email := "test@example.com"
	password := ""

	cfg := config.GetInstance().Config
	authClient := app.NewAuthClient(app.GetAuthConnection(log, cfg.Client))

	_, err := authClient.Login(context.Background(), email, password)
	assert.Error(t, err, "Expected error for empty password")
}
func TestLoginWithValidEmailAndPasswordButFailedToSaveToken(t *testing.T) {
	email := "test@example.com"
	password := "testpassword"

	mockAuthClient := &MockAuthClient{}
	mockAuthClient.On("Login", context.Background(), email, password).Return("test_token", nil)

	token, err := mockAuthClient.Login(context.Background(), email, password)
	assert.Nil(t, err, "Expected no error, but got %v", err)
	assert.Equal(t, "test_token", token, "Expected token 'test_token', but got %s", token)

	mockTokenStorage := &MockTokenStorage{}
	mockTokenStorage.On("Save", token).Return(errors.New("failed to save token"))

	err = mockTokenStorage.Save(token)
	assert.Error(t, err, "Expected error for failed to save token")

	mockAuthClient.AssertExpectations(t)
	mockTokenStorage.AssertExpectations(t)
}

func TestLoginWithValidEmailAndPasswordButAuthClientReturnsError(t *testing.T) {
	email := "test@example.com"
	password := "testpassword"

	mockAuthClient := &MockAuthClient{}
	expectedError := errors.New("failed to login")
	mockAuthClient.On("Login", context.Background(), email, password).Return("", expectedError)

	_, err := mockAuthClient.Login(context.Background(), email, password)
	assert.Error(t, err, "Expected error for failed login")
	assert.Equal(t, expectedError, err, "Expected error 'failed to login', but got %v", err)

	mockAuthClient.AssertExpectations(t)
}
func TestLoginWithValidEmailAndPasswordButTokenStorageReturnsError(t *testing.T) {
	email := "test@example.com"
	password := "testpassword"

	mockAuthClient := &MockAuthClient{}
	mockAuthClient.On("Login", context.Background(), email, password).Return("test_token", nil)

	token, err := mockAuthClient.Login(context.Background(), email, password)
	assert.Nil(t, err, "Expected no error, but got %v", err)
	assert.Equal(t, "test_token", token, "Expected token 'test_token', but got %s", token)

	mockTokenStorage := &MockTokenStorage{}
	expectedError := errors.New("failed to save token")
	mockTokenStorage.On("Save", token).Return(expectedError)

	err = mockTokenStorage.Save(token)
	assert.Error(t, err, "Expected error for failed to save token")
	assert.Equal(t, expectedError, err, "Expected error 'failed to save token', but got %v", err)

	mockAuthClient.AssertExpectations(t)
	mockTokenStorage.AssertExpectations(t)
}
