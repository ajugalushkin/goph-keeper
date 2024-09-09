package services

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

type JWTManager struct {
	log           *slog.Logger
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(log *slog.Logger, secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		log,
		secretKey,
		tokenDuration,
	}
}

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

func (manager *JWTManager) Verify(accessToken string) (bool, int64, error) {
	const op = "JWTManager.Verify"
	log := manager.log.With("op", op)

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

	if token == nil {

	}

	if err != nil {
		return false, 0, fmt.Errorf("invalid token: %w", err)
	}

	if token == nil {
		return false, 0, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		log.Debug("Failed to verify token",
			"error", err,
			"key", manager.secretKey,
			"valid", token.Valid)
		return false, 0, fmt.Errorf("invalid token: %w", err)
	}

	var userID int64
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID, err = strconv.ParseInt(fmt.Sprint(claims["uid"]), 10, 64)
		if err != nil {
			return false, 0, fmt.Errorf("invalid user ID: %w", err)
		}
	}

	log.Info("User ID", "uid", userID)
	return true, userID, nil
}
