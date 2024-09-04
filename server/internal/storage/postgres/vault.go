package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage"
)

type VaultStorage struct {
	db *sql.DB
}

func NewVaultStorage(storagePath string) (*VaultStorage, error) {
	const op = "storage.postgres.NewUserStorage"
	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &VaultStorage{db: db}, nil
}

func (v *VaultStorage) Create(ctx context.Context, item *models.Item) (*models.Item, error) {
	row := v.db.QueryRowContext(
		ctx,
		`INSERT INTO vaults (name, content, owner_id)
                   VALUES($1, $2, $3)
                   ON CONFLICT DO NOTHING RETURNING version`,
		item.Name, item.Content, item.OwnerID,
	)
	err := row.Scan(&item.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return item, storage.ErrItemConflict
	}
	return item, err
}

func (v *VaultStorage) Get(ctx context.Context, name string, userID int64) (*models.Item, error) {
	row := v.db.QueryRowContext(
		ctx,
		`SELECT content, version FROM vaults WHERE name = ($1) AND owner_id = ($2)`,
		name, userID,
	)
	secret := &models.Item{
		Name:    name,
		OwnerID: userID,
	}
	err := row.Scan(&secret.Content, &secret.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrItemNotFound
	}
	return secret, err
}

func (v *VaultStorage) List(ctx context.Context) []models.Item {
	//TODO implement me
	panic("implement me")
}
