package auth

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/ajugalushkin/goph-keeper/internal/dto/models"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	tokenTTl    string
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
}

// New returns a new instance of the Auth service.
func New(
	log *log.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL string,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		tokenTTl:    tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (token string, err error) {
	panic("implement me")
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (userID int64, err error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(slog.String("operation", op))

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("error", err.Error()))
		return 0, fmt.Errorf(`%s:%w`, op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		return 0, fmt.Errorf(`%s:%w`, op, err)
	}

	log.Info("user registered successfully")
	return id, nil
}
