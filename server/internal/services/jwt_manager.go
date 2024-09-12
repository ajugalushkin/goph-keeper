package services

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

//go:generate mockery --name UserProvider
type JWTManager struct {
	log           *slog.Logger
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTManager creates a new JWTManager instance with the provided logger, secret key, and token duration.
// The JWTManager is responsible for generating and verifying JWT tokens for user authentication.
//
// log: The logger instance to be used for logging.
// secretKey: The secret key used for signing and verifying JWT tokens.
// tokenDuration: The duration for which the generated JWT tokens should be valid.
//
// Returns a pointer to a new JWTManager instance.
func NewJWTManager(
    log *slog.Logger,
    secretKey string,
    tokenDuration time.Duration,
) *JWTManager {
    return &JWTManager{
        log,
        secretKey,
        tokenDuration,
    }
}

// NewToken generates a new JWT token for the provided user.
// The token is signed with the secret key and has an expiration time specified by the tokenDuration.
//
// user: The user for whom the JWT token should be generated.
//
// Returns:
// - A string representing the generated JWT token.
// - An error if there was an issue generating the token.
func (manager *JWTManager) NewToken(user models.User) (string, error) {
    const op = "JWTManager.NewToken"
    log := manager.log.With("op", op)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256,
        jwt.MapClaims{
            "uid":   user.ID,
            "email": user.Email,
            "exp":   time.Now().Add(manager.tokenDuration).Unix(),
        })

    tokenString, err := token.SignedString([]byte(manager.secretKey))
    if err != nil {
        log.Debug("Failed to sign token", "error", err)
        return "", err
    }

    log.Debug("JWT token generated", "token", tokenString, "secretKey", manager.secretKey)
    return tokenString, nil
}

// Verify verifies the provided JWT access token and returns the validity status, user ID, and any error encountered.
//
// accessToken: The JWT access token to be verified.
//
// Returns:
// - A boolean indicating whether the token is valid (true) or not (false).
// - An int64 representing the user ID extracted from the token.
// - An error if there was an issue verifying the token or extracting the user ID.
func (manager *JWTManager) Verify(accessToken string) (bool, int64, error) {
    const op = "JWTManager.Verify"
    log := manager.log.With("op", op)

    // Parse the JWT token using the provided secret key
    token, err := jwt.Parse(
        accessToken,
        func(token *jwt.Token) (interface{}, error) {
            _, ok := token.Method.(*jwt.SigningMethodHMAC)
            if !ok {
                log.Debug("Unexpected signing method")
                return nil, fmt.Errorf("unexpected token signing method")
            }

            return []byte(manager.secretKey), nil
        },
    )

    if err != nil {
        return false, 0, fmt.Errorf("invalid token: %w", err)
    }

    // Check if the token is nil
    if token == nil {
        return false, 0, fmt.Errorf("invalid token: %w", err)
    }

    // Check if the token is valid
    if !token.Valid {
        log.Debug("Failed to verify token",
            "error", err,
            "key", manager.secretKey,
            "valid", token.Valid)
        return false, 0, fmt.Errorf("invalid token: %w", err)
    }

    // Extract the user ID from the token claims
    var userID int64
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        userID, err = strconv.ParseInt(fmt.Sprint(claims["uid"]), 10, 64)
        if err != nil {
            return false, 0, fmt.Errorf("invalid user ID: %w", err)
        }
    }

    // Log the user ID
    log.Info("User ID", "uid", userID)
    return true, userID, nil
}
