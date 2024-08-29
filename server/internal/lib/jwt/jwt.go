package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"uid":   user.ID,
			"email": user.Email,
			"exp":   time.Now().Add(duration).Unix(),
		})

	tokenSecret := os.Getenv("TOKEN_SECRET")
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
