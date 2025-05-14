package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	jwtManager  TokenManager
}

//go:generate mockery --name UserSaver
type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

//go:generate mockery --name UserProvider
type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
}

var (
	ErrInvalidCredentials = errors.New("invalid creds")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

// NewAuthService returns a new instance of the Auth service.
// NewAuthService creates a new instance of the Auth service.
// The Auth service provides methods for user authentication and registration.
//
// Parameters:
// - log: A pointer to a slog.Logger instance for logging.
// - userSaver: An implementation of the UserSaver interface for saving user data.
// - userProvider: An implementation of the UserProvider interface for retrieving user data.
// - jwtManager: An implementation of the TokenManager interface for managing JWT tokens.
//
// Returns:
// - A pointer to a new instance of the Auth service.
func NewAuthService(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	jwtManager TokenManager,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		jwtManager:  jwtManager,
	}
}

// Login attempts to authenticate a user with the provided email and password.
// It retrieves the user from the user provider using the provided email.
// If the user is not found, it returns an error with the ErrInvalidCredentials code.
// If the user is found, it compares the provided password with the stored password hash.
// If the passwords do not match, it returns an error with the ErrInvalidCredentials code.
// If the passwords match, it generates a new JWT token_cache for the user using the JWT manager.
// If the token_cache generation fails, it returns an error.
// If the user is successfully authenticated, it logs the successful login and returns the generated token_cache.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "keeper.Login"
	log := a.log.With(slog.String("operation", op))

	log.Info("attempting to login user")
	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid creds", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := a.jwtManager.NewToken(user)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("Created token_cache", slog.String("token_cache", token))

	return token, nil
}

// RegisterNewUser registers a new user with the provided email and password.
// It generates a password hash using bcrypt and saves the user to the storage.
//
// Parameters:
// - ctx: A context.Context for cancellation and deadline support.
// - email: A string representing the user's email address.
// - password: A string representing the user's password.
//
// Returns:
//   - userID: An int64 representing the unique identifier of the registered user.
//   - err: An error if the registration process fails. If the user already exists,
//     it returns ErrUserExists. If any other error occurs during the registration process,
//     it returns the corresponding error.
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (userID int64, err error) {
	const op = "keeper.RegisterNewUser"
	log := a.log.With(slog.String("operation", op))

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("error", err.Error()))
		return 0, fmt.Errorf(`%s:%w`, op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", slog.String("error", err.Error()))
			return 0, fmt.Errorf(`%s: %w`, op, ErrUserExists)
		}

		a.log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, fmt.Errorf(`%s:%w`, op, err)
	}

	log.Info("user registered successfully")
	return id, nil
}
