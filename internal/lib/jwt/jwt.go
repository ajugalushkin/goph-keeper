package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ajugalushkin/goph-keeper/internal/dto/models"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenKey := os.Getenv("TOKEN_SECRET")
	tokenString, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
