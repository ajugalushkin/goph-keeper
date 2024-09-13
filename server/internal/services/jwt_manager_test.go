package services

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

// Creates a JWTManager instance with valid logger, secret key, and token_cache duration
func TestNewJWTManagerWithValidInputs(t *testing.T) {

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	secretKey := "mysecretkey"
	tokenDuration := time.Hour

	manager := NewJWTManager(log, secretKey, tokenDuration)

	if manager.log != log {
		t.Errorf("Expected logger %v, got %v", log, manager.log)
	}
	if manager.secretKey != secretKey {
		t.Errorf("Expected secretKey %v, got %v", secretKey, manager.secretKey)
	}
	if manager.tokenDuration != tokenDuration {
		t.Errorf("Expected tokenDuration %v, got %v", tokenDuration, manager.tokenDuration)
	}
}

// Handles empty secret key gracefully
func TestNewJWTManagerWithEmptySecretKey(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	secretKey := ""
	tokenDuration := time.Hour

	manager := NewJWTManager(log, secretKey, tokenDuration)

	if manager.secretKey != secretKey {
		t.Errorf("Expected empty secretKey, got %v", manager.secretKey)
	}
}

// Generates a valid JWT token_cache for a given user
func TestGeneratesValidJWTToken(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	secretKey := "test_secret_key"
	tokenDuration := time.Hour
	manager := NewJWTManager(log, secretKey, tokenDuration)

	user := models.User{
		ID:    1,
		Email: "test@example.com",
	}

	token, err := manager.NewToken(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected a valid token_cache, got an empty string")
	}
}

// Verify valid token_cache with correct signing method
func TestVerifyValidToken(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	manager := NewJWTManager(log, "secret", time.Hour)

	user := models.User{ID: 1, Email: "test@example.com"}
	token, err := manager.NewToken(user)
	if err != nil {
		t.Fatalf("Failed to create token_cache: %v", err)
	}

	valid, userID, err := manager.Verify(token)
	if err != nil {
		t.Fatalf("Failed to verify token_cache: %v", err)
	}

	if !valid {
		t.Errorf("Expected token_cache to be valid")
	}

	if userID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, userID)
	}
}

// Handle invalid token_cache format
func TestVerifyInvalidTokenFormat(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	manager := NewJWTManager(log, "secret", time.Hour)

	invalidToken := "invalid.token_cache.format"

	valid, userID, err := manager.Verify(invalidToken)

	if valid {
		t.Errorf("Expected token_cache to be invalid")
	}

	if userID != 0 {
		t.Errorf("Expected user ID to be 0, got %d", userID)
	}

	if err == nil {
		t.Errorf("Expected an error for invalid token_cache format")
	}
}
