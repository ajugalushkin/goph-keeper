package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage"
)

type UserStorage struct {
	db *sql.DB
}

// NewUserStorage creates a new instance of UserStorage using PostgreSQL database.
// It takes a storagePath parameter, which is the connection string for the PostgreSQL database.
// The function opens a connection to the database, pings it to ensure it's working,
// and returns a new UserStorage instance or an error if any occurs.
//
// storagePath: The connection string for the PostgreSQL database.
//
// Returns:
// - A pointer to a new UserStorage instance if successful.
// - An error if the database connection fails or cannot be pinged.
func NewUserStorage(storagePath string) (*UserStorage, error) {
	const op = "storage.postgres.NewUserStorage"
	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &UserStorage{db: db}, nil
}

// SaveUser inserts a new user into the database with the provided email and password hash.
// It returns the ID of the newly created user and an error if any occurs.
// If a user with the same email already exists in the database, it returns storage.ErrUserExists.
//
// ctx: The context for the operation.
// email: The email of the new user.
// passHash: The password hash of the new user.
//
// Returns:
// - uid: The ID of the newly created user.
// - err: An error if any occurs during the operation.
func (s *UserStorage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (uid int64, err error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, password_hash) VALUES ($1, $2)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	lastUser, err := s.User(ctx, email)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return lastUser.ID, nil
}

// User retrieves a user from the database by their email.
// It returns the user's information and an error if any occurs.
// If no user is found with the given email, it returns storage.ErrUserNotFound.
//
// ctx: The context for the operation.
// email: The email of the user to retrieve.
//
// Returns:
// - user: The user's information if found.
// - err: An error if any occurs during the operation.
func (s *UserStorage) User(
	ctx context.Context,
	email string,
) (models.User, error) {
	const op = "storage.postgres.User"

	// Prepare a SQL statement to select user information by email.
	stmt, err := s.db.Prepare("SELECT id, email, password_hash FROM users WHERE email = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	// Execute the prepared statement with the given email.
	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		// If no rows were returned, return storage.ErrUserNotFound.
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		// If any other error occurred, return it.
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	// Return the retrieved user information.
	return user, nil
}
