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

func (v *VaultStorage) List(ctx context.Context) []models.Item {
	//TODO implement me
	panic("implement me")
}
