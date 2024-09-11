package storage

import (
	"context"
	"errors"
	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

var (
	ErrUserExists   = errors.New("already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrUserConflict = errors.New("user conflict")
	ErrItemConflict = errors.New("item conflict")
	ErrItemNotFound = errors.New("item not found")
)

//go:generate mockery --name UserStorage
type UserStorage interface {
	User(ctx context.Context, email string) (user models.User, err error)
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}
