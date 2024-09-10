package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage"
)

// Handles invalid storage path gracefully
func TestNewUserStorage_InvalidPath(t *testing.T) {
	storagePath := "invalid_path"
	userStorage, err := NewUserStorage(storagePath)

	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if userStorage != nil {
		t.Fatalf("expected userStorage to be nil, got %v", userStorage)
	}
}

// Save a new user with a unique email and password hash
func TestSaveUserWithUniqueEmail(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStorage := &UserStorage{db: db}
	email := "unique@example.com"
	passHash := []byte("hashedpassword")

	mock.ExpectPrepare("INSERT INTO users\\(email, password_hash\\) VALUES \\(\\$1, \\$2\\)").ExpectExec().WithArgs(email, passHash).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("SELECT id, email, password_hash FROM users WHERE email = \\$1").ExpectQuery().WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).AddRow(1, email, passHash))

	uid, err := userStorage.SaveUser(ctx, email, passHash)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if uid != 1 {
		t.Errorf("expected user ID to be 1, got %d", uid)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Attempt to save a user with an email that already exists
func TestSaveUserWithExistingEmail(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStorage := &UserStorage{db: db}
	email := "existing@example.com"
	passHash := []byte("hashedpassword")

	mock.ExpectPrepare("INSERT INTO users\\(email, password_hash\\) VALUES \\(\\$1, \\$2\\)").ExpectExec().WithArgs(email, passHash).WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

	_, err = userStorage.SaveUser(ctx, email, passHash)
	if !errors.Is(err, storage.ErrUserExists) {
		t.Errorf("expected error to be %v, got %v", storage.ErrUserExists, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Retrieves user details successfully when the email exists in the database
func TestUserRetrievesUserDetailsSuccessfully(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	email := "test@example.com"
	user := models.User{
		ID:           1,
		Email:        email,
		PasswordHash: []byte("hashedpassword"),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash"}).
		AddRow(user.ID, user.Email, user.PasswordHash)

	mock.ExpectPrepare("SELECT id, email, password_hash FROM users WHERE email = \\$1").
		ExpectQuery().
		WithArgs(email).
		WillReturnRows(rows)

	userStorage := &UserStorage{db: db}
	_, err = userStorage.User(ctx, email)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

// Returns ErrUserNotFound when the email does not exist in the database
func TestUserReturnsErrUserNotFound(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	email := "nonexistent@example.com"

	mock.ExpectPrepare("SELECT id, email, password_hash FROM users WHERE email = \\$1").
		ExpectQuery().
		WithArgs(email).
		WillReturnError(pgx.ErrNoRows)

	userStorage := &UserStorage{db: db}
	_, err = userStorage.User(ctx, email)
	if err == nil {
		t.Error("expected an error but got none")
	}

	if !errors.Is(err, storage.ErrUserNotFound) {
		t.Errorf("expected %v, got %v", storage.ErrUserNotFound, err)
	}
}
