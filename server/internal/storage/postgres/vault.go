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

func (v *VaultStorage) Update(ctx context.Context, item *models.Item) (*models.Item, error) {
	SQLQuery := `
        UPDATE secrets
        SET version = uuid_generate_v4(), content = ($1)
        WHERE owner_id = ($2) AND name = ($3)
        RETURNING version`

	row := v.db.QueryRowContext(ctx, SQLQuery, item.Content, item.OwnerID, item.Name)
	err := row.Scan(&item.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrItemNotFound
		}
		return nil, err
	}

	return item, nil
}

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
