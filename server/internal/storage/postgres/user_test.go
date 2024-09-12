package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
func TestSaveUserWithEmailExceedingMaxLength(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStorage := &UserStorage{db: db}
	email := strings.Repeat("a", 256) + "@example.com" // 256 characters
	passHash := []byte("hashedpassword")

	mock.ExpectPrepare("INSERT INTO users\\(email, password_hash\\) VALUES \\(\\$1, \\$2\\)").
		ExpectExec().
		WithArgs(email, passHash).
		WillReturnError(fmt.Errorf("value too long for type character varying(255)"))

	_, err = userStorage.SaveUser(ctx, email, passHash)
	if err == nil {
		t.Error("expected an error but got none")
	} else {
		expectedError := "value too long for type character varying(255)"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected error to contain '%s', got '%s'", expectedError, err.Error())
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestSaveUserWithPasswordHashExceedingMaxLength(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStorage := &UserStorage{db: db}
	email := "test@example.com"
	passHash := make([]byte, 513) // 513 bytes, which exceeds the maximum length

	mock.ExpectPrepare("INSERT INTO users\\(email, password_hash\\) VALUES \\(\\$1, \\$2\\)").
		ExpectExec().
		WithArgs(email, passHash).
		WillReturnError(fmt.Errorf("value too long for type bytea"))

	_, err = userStorage.SaveUser(ctx, email, passHash)
	if err == nil {
		t.Error("expected an error but got none")
	} else {
		expectedError := "value too long for type bytea"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected error to contain '%s', got '%s'", expectedError, err.Error())
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestRetrieveUserDetailsWithSQLInjectionAttempt(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStorage := &UserStorage{db: db}
	email := "test@example.com' OR 1=1 --" // SQL injection attempt

	mock.ExpectPrepare("SELECT id, email, password_hash FROM users WHERE email = \\$1").
		ExpectQuery().
		WithArgs(email).
		WillReturnError(pgx.ErrNoRows) // Simulate no user found

	_, err = userStorage.User(ctx, email)
	if err == nil {
		t.Error("expected an error but got none")
	}

	if !errors.Is(err, storage.ErrUserNotFound) {
		t.Errorf("expected %v, got %v", storage.ErrUserNotFound, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveUserWithInvalidPasswordHashCharacters(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStorage := &UserStorage{db: db}
	email := "test@example.com"
	passHash := []byte("invalid_characters") // Invalid characters

	mock.ExpectPrepare("INSERT INTO users\\(email, password_hash\\) VALUES \\(\\$1, \\$2\\)").
		ExpectExec().
		WithArgs(email, passHash).
		WillReturnError(fmt.Errorf("invalid characters in password hash"))

	_, err = userStorage.SaveUser(ctx, email, passHash)
	if err == nil {
		t.Error("expected an error but got none")
	} else {
		expectedError := "invalid characters in password hash"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected error to contain '%s', got '%s'", expectedError, err.Error())
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveUser_EmptyEmail(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStorage := &UserStorage{db: db}
	email := ""
	passHash := []byte("hashedpassword")

	mock.ExpectPrepare("INSERT INTO users\\(email, password_hash\\) VALUES \\(\\$1, \\$2\\)").
		ExpectExec().
		WithArgs(email, passHash).
		WillReturnError(fmt.Errorf("empty email"))

	_, err = userStorage.SaveUser(ctx, email, passHash)
	if err == nil {
		t.Error("expected an error but got none")
	} else {
		expectedError := "empty email"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected error to contain '%s', got '%s'", expectedError, err.Error())
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveUserWithEmptyPasswordHash(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStorage := &UserStorage{db: db}
	email := "test@example.com"
	passHash := []byte{} // Empty password hash

	mock.ExpectPrepare("INSERT INTO users\\(email, password_hash\\) VALUES \\(\\$1, \\$2\\)").
		ExpectExec().
		WithArgs(email, passHash).
		WillReturnError(fmt.Errorf("password hash cannot be empty"))

	_, err = userStorage.SaveUser(ctx, email, passHash)
	if err == nil {
		t.Error("expected an error but got none")
	} else {
		expectedError := "password hash cannot be empty"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected error to contain '%s', got '%s'", expectedError, err.Error())
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveUserWithInvalidEmailCharacters(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStorage := &UserStorage{db: db}
	email := "test@example.com!#" // Invalid characters
	passHash := []byte("hashedpassword")

	mock.ExpectPrepare("INSERT INTO users\\(email, password_hash\\) VALUES \\(\\$1, \\$2\\)").
		ExpectExec().
		WithArgs(email, passHash).
		WillReturnError(fmt.Errorf("invalid characters in email"))

	_, err = userStorage.SaveUser(ctx, email, passHash)
	if err == nil {
		t.Error("expected an error but got none")
	} else {
		expectedError := "invalid characters in email"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected error to contain '%s', got '%s'", expectedError, err.Error())
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveUserWithSQLInjectionAttempt(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStorage := &UserStorage{db: db}
	email := "test@example.com' OR 1=1 --" // SQL injection attempt
	passHash := []byte("hashedpassword")

	mock.ExpectPrepare("INSERT INTO users\\(email, password_hash\\) VALUES \\(\\$1, \\$2\\)").
		ExpectExec().
		WithArgs(email, passHash).
		WillReturnError(fmt.Errorf("unique constraint violation"))

	_, err = userStorage.SaveUser(ctx, email, passHash)
	if err == nil {
		t.Error("expected an error but got none")
	} else {
		expectedError := "unique constraint violation"
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected error to contain '%s', got '%s'", expectedError, err.Error())
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
