package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ajugalushkin/goph-keeper/internal/dto/models"
)

func NewToken(user models.User, duration time.Duration, tokenSecret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
