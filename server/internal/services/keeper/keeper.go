package keeper

import (
	"context"
	"log/slog"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

type Keeper struct {
	log         *slog.Logger
	itmSaver    ItemSaver
	itmProvider ItemProvider
}

type ItemProvider interface {
	List(ctx context.Context) []models.Item
}

type ItemSaver interface {
	Create(ctx context.Context, item *models.Item) (*models.Item, error)
}

func New(log *slog.Logger, provider ItemProvider, saver ItemSaver) *Keeper {
	return &Keeper{
		log:         log,
		itmSaver:    saver,
		itmProvider: provider,
	}
}

func (k Keeper) CreateItem(ctx context.Context, item *models.Item) (*models.Item, error) {
	newItem, err := k.itmSaver.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return newItem, nil
}

func (k Keeper) ListItem(ctx context.Context, since int64) (list *models.ListItem, err error) {
	k.itmProvider.List(ctx)
	return nil, nil
}
