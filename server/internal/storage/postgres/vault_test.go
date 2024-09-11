package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage"
)

// Handles the case where the item already exists and returns an item conflict error
func TestCreateItemConflict(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	vaultStorage := &VaultStorage{db: db}
	item := &models.Item{
		Name:    "test_item",
		Content: []byte("test_content"),
		OwnerID: 1,
		FileID:  "file_123",
	}

	mock.ExpectQuery(`INSERT INTO vaults`).
		WithArgs(item.Name, item.Content, item.OwnerID, item.FileID).
		WillReturnError(sql.ErrNoRows)

	_, err = vaultStorage.Create(ctx, item)
	if !errors.Is(err, storage.ErrItemConflict) {
		t.Errorf("expected error %v, got %v", storage.ErrItemConflict, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Successfully deletes an item when provided with valid name and owner_id
func TestDeleteItemSuccess(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	item := &models.Item{
		Name:    "test_item",
		OwnerID: 1,
	}

	mock.ExpectExec(`DELETE FROM vaults WHERE name = \(\$1\) AND owner_id = \(\$2\)`).
		WithArgs(item.Name, item.OwnerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	vaultStorage := &VaultStorage{db: db}
	err = vaultStorage.Delete(ctx, item)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Returns error if the database connection is lost during execution
func TestDeleteItemDBConnectionLost(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	item := &models.Item{
		Name:    "test_item",
		OwnerID: 1,
	}

	mock.ExpectExec(`DELETE FROM vaults WHERE name = \(\$1\) AND owner_id = \(\$2\)`).
		WithArgs(item.Name, item.OwnerID).
		WillReturnError(fmt.Errorf("db connection lost"))

	vaultStorage := &VaultStorage{db: db}
	err = vaultStorage.Delete(ctx, item)
	if err == nil {
		t.Error("expected an error but got none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Item does not exist in the database
func TestGetItemNotFound(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	vaultStorage := &VaultStorage{db: db}
	itemName := "nonExistentItem"
	userID := int64(1)

	mock.ExpectQuery(`SELECT content, version, file_id FROM vaults WHERE name = \(\$1\) AND owner_id = \(\$2\)`).
		WithArgs(itemName, userID).
		WillReturnError(sql.ErrNoRows)

	item, err := vaultStorage.Get(ctx, itemName, userID)
	if err != storage.ErrItemNotFound {
		t.Errorf("expected error %v, got %v", storage.ErrItemNotFound, err)
	}

	if item != nil {
		t.Errorf("expected nil item, got %v", item)
	}
}

// Handles database connection errors gracefully
func TestListHandlesDBConnectionError(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT name, version, content FROM vaults WHERE owner_id = \(\$1\)`).
		WithArgs(userID).
		WillReturnError(fmt.Errorf("db connection error"))

	vaultStorage := &VaultStorage{db: db}
	_, err = vaultStorage.List(ctx, userID)
	if err == nil {
		t.Error("expected an error but got none")
	}

	if err.Error() != "db connection error" {
		t.Errorf("expected 'db connection error', got %v", err)
	}
}
