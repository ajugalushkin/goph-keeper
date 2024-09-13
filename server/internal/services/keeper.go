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

//go:generate mockery --name ItemProvider
type ItemProvider interface {
	Get(ctx context.Context, name string, userID int64) (*models.Item, error)
	List(ctx context.Context, userID int64) ([]*models.Item, error)
}

//go:generate mockery --name ItemSaver
type ItemSaver interface {
	Create(ctx context.Context, item *models.Item) (*models.Item, error)
	Update(ctx context.Context, item *models.Item) (*models.Item, error)
	Delete(ctx context.Context, item *models.Item) error
}

//go:generate mockery --name ObjectSaver
type ObjectSaver interface {
	Create(ctx context.Context, file *models.File) (string, error)
	Delete(ctx context.Context, objectID string) error
}

//go:generate mockery --name ObjectProvider
type ObjectProvider interface {
	Get(ctx context.Context, objectID string) (*models.File, error)
}

// NewKeeperService creates a new instance of the Keeper service.
// The Keeper service provides methods for managing items and files.
//
// Parameters:
// - log: A pointer to a slog.Logger instance for logging.
// - provider: An implementation of the ItemProvider interface for retrieving item data.
// - saver: An implementation of the ItemSaver interface for saving item data.
// - objectSaver: An implementation of the ObjectSaver interface for saving file objects.
// - objectProvider: An implementation of the ObjectProvider interface for retrieving file objects.
//
// Returns:
// - A pointer to a new instance of the Keeper service.
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

// CreateItem creates a new item in the keeper.
//
// The function accepts a context and an item as parameters. The context is used to control the
// execution deadline, cancelation, and other request-scoped values. The item parameter represents
// the data to be stored in the keeper.
//
// The function interacts with the ItemSaver interface to save the new item. If the save operation
// is successful, the function returns the newly created item and a nil error. If an error occurs
// during the save operation, the function logs the error, returns nil for the item, and the error.
//
// The function also logs the operation using the provided logger (k.log).
//
// Parameters:
// - ctx: A context.Context used to control the execution of the function.
// - item: A pointer to a models.Item representing the data to be stored in the keeper.
//
// Returns:
// - A pointer to a models.Item representing the newly created item.
// - An error if the save operation fails.
func (k *Keeper) CreateItem(
	ctx context.Context,
	item *models.Item,
) (*models.Item, error) {
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

// CreateFile creates a new file in the keeper and associates it with an item.
//
// The function accepts a context and a file as parameters. The context is used to control the
// execution deadline, cancelation, and other request-scoped values. The file parameter represents
// the data to be stored in the keeper.
//
// The function interacts with the ObjectSaver and ItemSaver interfaces to save the file and item data,
// respectively. If the save operations are successful, the function returns the version of the newly
// created item and a nil error. If an error occurs during the save operations, the function logs the
// error, returns an empty string for the version, and the error.
//
// Parameters:
// - ctx: A context.Context used to control the execution of the function.
// - file: A pointer to a models.File representing the data to be stored in the keeper.
//
// Returns:
// - A string representing the version of the newly created item.
// - An error if the save operations fail.
func (k *Keeper) CreateFile(
	ctx context.Context,
	file *models.File,
) (string, error) {
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

// UpdateItem updates an existing item in the keeper.
//
// The function accepts a context and an item as parameters. The context is used to control the
// execution deadline, cancelation, and other request-scoped values. The item parameter represents
// the data to be updated in the keeper.
//
// The function interacts with the ItemSaver interface to update the item. If the update operation
// is successful, the function returns the updated item and a nil error. If an error occurs during
// the update operation, the function logs the error, returns nil for the item, and the error.
//
// The function also logs the operation using the provided logger (k.log).
//
// Parameters:
// - ctx: A context.Context used to control the execution of the function.
// - item: A pointer to a models.Item representing the data to be updated in the keeper.
//
// Returns:
// - A pointer to a models.Item representing the updated item.
// - An error if the update operation fails.
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

// DeleteItem deletes an existing item in the keeper along with its associated file, if any.
//
// The function accepts a context and an item as parameters. The context is used to control the
// execution deadline, cancelation, and other request-scoped values. The item parameter represents
// the data to be deleted from the keeper.
//
// The function interacts with the ItemProvider and ItemSaver interfaces to retrieve and del the item.
// If the item has an associated file (identified by the FileID field), the function also interacts with
// the ObjectSaver interface to del the file.
//
// If any of the deletion operations fail, the function logs the error and returns the error.
// Otherwise, the function logs a success message and returns nil.
//
// Parameters:
// - ctx: A context.Context used to control the execution of the function.
// - item: A pointer to a models.Item representing the data to be deleted from the keeper.
//
// Returns:
// - An error if any of the deletion operations fail.
func (k *Keeper) DeleteItem(
	ctx context.Context,
	item *models.Item,
) error {
	const op = "services.keeper.deleteItem"
	log := k.log.With("op", op)

	item, err := k.itmProvider.Get(ctx, item.Name, item.OwnerID)
	if err != nil {
		log.Debug("failed get item", slog.String("error", err.Error()))
		return err
	}

	err = k.itmSaver.Delete(ctx, item)
	if err != nil {
		k.log.Debug("Failed to del item", slog.String("error", err.Error()))
		return err
	}

	if item.FileID != "" {
		err := k.ObjSaver.Delete(ctx, item.FileID)
		if err != nil {
			log.Debug("Failed to del file", slog.String("error", err.Error()))
			return err
		}
	}

	k.log.Debug("Successfully deleted item")
	return nil
}

// GetItem retrieves an item from the keeper by name and user ID.
//
// The function interacts with the ItemProvider interface to fetch the item data.
// If the item is found, it is returned along with a nil error. If an error occurs during the retrieval,
// the function logs the error, returns nil for the item, and the error.
//
// Parameters:
// - ctx: A context.Context used to control the execution of the function.
// - name: A string representing the name of the item to retrieve.
// - userID: An int64 representing the user ID associated with the item.
//
// Returns:
// - A pointer to a models.Item representing the retrieved item.
// - An error if the retrieval operation fails.
func (k *Keeper) GetItem(
	ctx context.Context,
	name string,
	userID int64,
) (*models.Item, error) {
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

// GetFile retrieves a file from the keeper by name and user ID.
// The function interacts with the ItemProvider and ObjectProvider interfaces to fetch the item and file data.
// If the item and file are found, they are combined into a models.File struct and returned along with a nil error.
// If an error occurs during the retrieval, the function logs the error, returns nil for the file, and the error.
//
// Parameters:
// - ctx: A context.Context used to control the execution of the function.
// - name: A string representing the name of the item to retrieve.
// - userID: An int64 representing the user ID associated with the item.
//
// Returns:
// - A pointer to a models.File representing the retrieved file.
// - An error if the retrieval operation fails.
func (k *Keeper) GetFile(
	ctx context.Context,
	name string,
	userID int64,
) (*models.File, error) {
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

// ListItems retrieves a list of items from the keeper for a given user ID.
//
// The function interacts with the ItemProvider interface to fetch the list of items associated with the provided user ID.
// If the retrieval operation is successful, the function returns a slice of models.Item pointers along with a nil error.
// If an error occurs during the retrieval, the function logs the error, returns nil for the list, and the error.
//
// Parameters:
// - ctx: A context.Context used to control the execution of the function.
// - userID: An int64 representing the user ID associated with the items to retrieve.
//
// Returns:
// - list: A slice of models.Item pointers representing the retrieved items.
// - err: An error if the retrieval operation fails. If nil, the operation was successful.
func (k *Keeper) ListItems(
	ctx context.Context,
	userID int64,
) (list []*models.Item, err error) {
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
