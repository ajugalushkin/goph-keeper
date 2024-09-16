package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage"
)

type VaultStorage struct {
	db *sql.DB
}

// NewVaultStorage creates a new instance of VaultStorage using a PostgreSQL database.
// It takes a storagePath parameter, which is the connection string for the PostgreSQL database.
// The function opens a connection to the database and returns a new VaultStorage instance or an error if the connection fails.
//
// storagePath: The connection string for the PostgreSQL database.
//
// Returns:
// - A pointer to a new VaultStorage instance if the connection is successful.
// - An error if the connection fails.
func NewVaultStorage(storagePath string) (*VaultStorage, error) {
	const op = "storage.postgres.NewUserStorage"
	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &VaultStorage{db: db}, nil
}

// Create inserts a new item into the vault storage.
// If an item with the same name and owner already exists, it returns the existing item and storage.ErrItemConflict.
//
// ctx: The context for the operation.
// item: The item to be inserted. The item's Name, Content, OwnerID, and FileID fields are required.
//
// Returns:
// - A pointer to the inserted item with its version field populated.
// - An error if the operation fails, which can be storage.ErrItemConflict if a conflict occurs.
func (v *VaultStorage) Create(
	ctx context.Context,
	item *models.Item,
) (*models.Item, error) {
	row := v.db.QueryRowContext(
		ctx,
		`INSERT INTO vaults (name, content, owner_id, file_id)
                   VALUES($1, $2, $3, $4)
                   ON CONFLICT DO NOTHING RETURNING version`,
		item.Name, item.Content, item.OwnerID, item.FileID,
	)
	err := row.Scan(&item.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return item, storage.ErrItemConflict
	}
	return item, err
}

// Update updates an existing item in the vault storage.
// It generates a new version for the item and updates its content.
// If the item with the specified name and owner does not exist, it returns storage.ErrItemNotFound.
//
// ctx: The context for the operation.
// item: The item to be updated. The item's Name, Content, and OwnerID fields are required.
//
// Returns:
// - A pointer to the updated item with its version field populated.
// - An error if the operation fails, which can be storage.ErrItemNotFound if the item does not exist.
func (v *VaultStorage) Update(ctx context.Context, item *models.Item) (*models.Item, error) {
	SQLQuery := `
        UPDATE vaults
        SET version = ($1), content = ($2)
        WHERE owner_id = ($3) AND name = ($4)
        RETURNING version`

	row := v.db.QueryRowContext(ctx, SQLQuery, uuid.New(), item.Content, item.OwnerID, item.Name)
	err := row.Scan(&item.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrItemNotFound
		}
		return nil, err
	}

	return item, nil
}

// Delete removes an existing item from the vault storage based on the provided name and owner ID.
// It executes a DELETE SQL query on the 'vaults' table with the specified name and owner ID.
//
// ctx: The context for the operation.
// item: The item to be deleted. The item's Name and OwnerID fields are required.
//
// Returns:
// - An error if the operation fails. If the item does not exist, it returns nil.
func (v *VaultStorage) Delete(
	ctx context.Context,
	item *models.Item,
) error {
	_, err := v.db.ExecContext(
		ctx,
		`DELETE FROM vaults WHERE name = ($1) AND owner_id = ($2)`,
		item.Name,
		item.OwnerID,
	)
	return err
}

// Get retrieves a single item from the vault storage based on the provided name and owner ID.
// It executes a SELECT SQL query on the 'vaults' table with the specified name and owner ID.
//
// ctx: The context for the operation. It is used to control the timeout and cancellation of the operation.
// name: The name of the item to retrieve.
// userID: The ID of the owner of the item.
//
// Returns:
// - A pointer to the retrieved item with its content, version, and file ID fields populated.
// - An error if the operation fails. If the item does not exist, it returns storage.ErrItemNotFound.
func (v *VaultStorage) Get(
	ctx context.Context,
	name string,
	userID int64,
) (*models.Item, error) {
	row := v.db.QueryRowContext(
		ctx,
		`SELECT content, version, file_id FROM vaults WHERE name = ($1) AND owner_id = ($2)`,
		name, userID,
	)
	secret := &models.Item{
		Name:    name,
		OwnerID: userID,
	}
	err := row.Scan(&secret.Content, &secret.Version, &secret.FileID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrItemNotFound
	}
	return secret, err
}

// List retrieves a list of items from the vault storage based on the provided owner ID.
// It executes a SELECT SQL query on the 'vaults' table with the specified owner ID.
//
// ctx: The context for the operation. It is used to control the timeout and cancellation of the operation.
// userID: The ID of the owner of the items to retrieve.
//
// Returns:
// - A slice of pointers to the retrieved items with their name, version, and content fields populated.
// - An error if the operation fails. If no items are found, it returns an empty slice and nil error.
func (v *VaultStorage) List(
	ctx context.Context,
	userID int64,
) ([]*models.Item, error) {
	rows, err := v.db.QueryContext(
		ctx, `SELECT name, version, content FROM vaults WHERE owner_id = ($1)`, userID)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	secrets := make([]*models.Item, 0)
	for rows.Next() {
		secret := &models.Item{
			OwnerID: userID,
		}
		if err = rows.Scan(&secret.Name, &secret.Version, &secret.Content); err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}
	return secrets, nil
}
