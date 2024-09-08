package services

import (
	"context"
	"log/slog"

	"github.com/minio/minio-go/v7"

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
}

type ObjectProvider interface {
	Get(ctx context.Context, userID int64, fileName string) (*minio.Object, error)
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
	return k.ObjSaver.Create(ctx, file)
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
	k.log.With("op", op)

	err := k.itmSaver.Delete(ctx, item)
	if err != nil {
		k.log.Debug("Failed to delete item", slog.String("error", err.Error()))
		return err
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

func (k *Keeper) GetFile(ctx context.Context, userID int64, fileName string) (*minio.Object, error) {
	return k.objProvider.Get(ctx, userID, fileName)
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
