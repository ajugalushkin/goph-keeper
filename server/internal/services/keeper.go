package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

type Keeper struct {
	log         *slog.Logger
	itmSaver    ItemSaver
	itmProvider ItemProvider
	ObjSaver    ObjectSaver
	objProvider ObjectProvider
}

type ItemProvider interface {
	Get(ctx context.Context, name string, userID int64) (*models.Item, error)
	List(ctx context.Context, userID int64) ([]*models.Item, error)
}

type ItemSaver interface {
	Create(ctx context.Context, item *models.Item) (*models.Item, error)
	Update(ctx context.Context, item *models.Item) (*models.Item, error)
	Delete(ctx context.Context, item *models.Item) error
}

type ObjectSaver interface {
	Create(ctx context.Context, file *models.File) (string, error)
	Delete(ctx context.Context, objectID string) error
}

type ObjectProvider interface {
	Get(ctx context.Context, objectID string) (*models.File, error)
}

func NewKeeperService(
	log *slog.Logger,
	provider ItemProvider,
	saver ItemSaver,
	objectSaver ObjectSaver,
	objectProvider ObjectProvider,
) *Keeper {
	return &Keeper{
		log:         log,
		itmSaver:    saver,
		itmProvider: provider,
		ObjSaver:    objectSaver,
		objProvider: objectProvider,
	}
}

func (k *Keeper) CreateItem(ctx context.Context, item *models.Item) (*models.Item, error) {
	const op = "services.keeper.createItem"
	k.log.With("op", op)

	newItem, err := k.itmSaver.Create(ctx, item)
	if err != nil {
		k.log.Debug("Failed to create item", slog.String("error", err.Error()))
		return nil, err
	}

	k.log.Debug("Successfully created item")
	return newItem, nil
}

func (k *Keeper) CreateFile(ctx context.Context, file *models.File) (string, error) {
	const op = "services.keeper.createFile"

	var err error
	file.Item.FileID, err = k.ObjSaver.Create(ctx, file)
	if err != nil || file.Item.FileID == "" {
		return "", fmt.Errorf("op: %s, failed to create file: %w", op, err)
	}

	item, err := k.itmSaver.Create(ctx, &file.Item)
	if err != nil {
		return "", fmt.Errorf("op: %s, failed to create info item for file: %w", op, err)
	}

	return item.Version.String(), nil
}

func (k *Keeper) UpdateItem(ctx context.Context, item *models.Item) (*models.Item, error) {
	const op = "services.keeper.updateItem"
	k.log.With("op", op)

	item, err := k.itmSaver.Update(ctx, item)
	if err != nil {
		k.log.Debug("Failed to update item", slog.String("error", err.Error()))
		return nil, err
	}

	k.log.Debug("Successfully update item")
	return item, nil
}

func (k *Keeper) DeleteItem(ctx context.Context, item *models.Item) error {
	const op = "services.keeper.deleteItem"
	log := k.log.With("op", op)

	item, err := k.itmProvider.Get(ctx, item.Name, item.OwnerID)
	if err != nil {
		log.Debug("failed get item", slog.String("error", err.Error()))
		return err
	}

	err = k.itmSaver.Delete(ctx, item)
	if err != nil {
		k.log.Debug("Failed to delete item", slog.String("error", err.Error()))
		return err
	}

	if item.FileID != "" {
		err := k.ObjSaver.Delete(ctx, item.FileID)
		if err != nil {
			log.Debug("Failed to delete file", slog.String("error", err.Error()))
			return err
		}
	}

	k.log.Debug("Successfully deleted item")
	return nil
}

func (k *Keeper) GetItem(ctx context.Context, name string, userID int64) (*models.Item, error) {
	const op = "services.keeper.getItem"
	k.log.With("op", op)

	item, err := k.itmProvider.Get(ctx, name, userID)
	if err != nil {
		k.log.Debug("Failed to get item", slog.String("error", err.Error()))
		return nil, err
	}

	k.log.Debug("Successfully get item")
	return item, nil
}

func (k *Keeper) GetFile(ctx context.Context, name string, userID int64) (*models.File, error) {
	const op = "services.keeper.getFile"

	item, err := k.itmProvider.Get(ctx, name, userID)
	if err != nil {
		return nil, fmt.Errorf("op: %s, failed to get info item for file: %w", op, err)
	}

	file, err := k.objProvider.Get(ctx, item.FileID)
	if err != nil {
		return nil, fmt.Errorf("op: %s, failed to get file: %w", op, err)
	}

	file.Item = *item

	return file, nil
}

func (k *Keeper) ListItems(ctx context.Context, userID int64) (list []*models.Item, err error) {
	const op = "services.keeper.listItem"
	k.log.With("op", op)

	list, err = k.itmProvider.List(ctx, userID)
	if err != nil {
		k.log.Debug("Failed to list items", slog.String("error", err.Error()))
		return nil, err
	}

	k.log.Debug("Successfully get list")
	return list, nil
}
