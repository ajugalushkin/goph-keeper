package services

import "github.com/ajugalushkin/goph-keeper/server/internal/dto/models"

//go:generate mockery --name TokenManager
type TokenManager interface {
	NewToken(user models.User) (string, error)
	Verify(accessToken string) (bool, int64, error)
}
