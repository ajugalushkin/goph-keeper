package services

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/services/mocks"
)

// Initializes Auth service with valid logger, userSaver, userProvider, and jwtManager
func TestNewAuthServiceInitialization(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	userSaver := new(mocks.UserSaver)
	userProvider := new(mocks.UserProvider)
	jwtManager := &JWTManager{
		log:           log,
		secretKey:     "secret",
		tokenDuration: time.Hour,
	}

	authService := NewAuthService(log, userSaver, userProvider, jwtManager)

	if authService.log != log {
		t.Errorf("expected logger to be %v, got %v", log, authService.log)
	}
	if authService.usrSaver != userSaver {
		t.Errorf("expected userSaver to be %v, got %v", userSaver, authService.usrSaver)
	}
	if authService.usrProvider != userProvider {
		t.Errorf("expected userProvider to be %v, got %v", userProvider, authService.usrProvider)
	}
	if authService.jwtManager != jwtManager {
		t.Errorf("expected jwtManager to be %v, got %v", jwtManager, authService.jwtManager)
	}
}

// Handles nil logger gracefully
func TestNewAuthServiceNilLogger(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	userSaver := new(mocks.UserSaver)
	userProvider := new(mocks.UserProvider)
	jwtManager := &JWTManager{
		log:           log,
		secretKey:     "secret",
		tokenDuration: time.Hour,
	}

	authService := NewAuthService(log, userSaver, userProvider, jwtManager)

	if authService.log != nil {
		t.Errorf("expected logger to be nil, got %v", authService.log)
	}
}

// Successful login returns a valid JWT token
func TestLoginSuccessful(t *testing.T) {
	ctx := context.Background()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{ID: 1, Email: email, PasswordHash: hashedPassword}

	mockUserProvider := new(mocks.UserProvider)
	mockUserProvider.On("User", ctx, email).Return(user, nil)

	mockJWTManager := new(mocks.TokenManager)
	mockJWTManager.On("NewToken", user).Return("valid-token", nil)

	auth := NewAuthService(log, nil, mockUserProvider, mockJWTManager)

	token, err := auth.Login(ctx, email, password)

	assert.NoError(t, err)
	assert.Equal(t, "valid-token", token)
}
