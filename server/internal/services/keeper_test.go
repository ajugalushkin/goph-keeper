package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/services/mocks"
)

// Initializes Keeper with valid dependencies
func TestNewKeeperService_ValidDependencies(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	provider := mocks.NewItemProvider(t)
	saver := mocks.NewItemSaver(t)
	objectSaver := mocks.NewObjectSaver(t)
	objectProvider := mocks.NewObjectProvider(t)

	keeper := NewKeeperService(log, provider, saver, objectSaver, objectProvider)

	if keeper == nil {
		t.Fatal("Expected keeper to be initialized, got nil")
	}
	if keeper.log != log {
		t.Errorf("Expected log to be %v, got %v", log, keeper.log)
	}
	if keeper.itmProvider != provider {
		t.Errorf("Expected itmProvider to be %v, got %v", provider, keeper.itmProvider)
	}
	if keeper.itmSaver != saver {
		t.Errorf("Expected itmSaver to be %v, got %v", saver, keeper.itmSaver)
	}
	if keeper.ObjSaver != objectSaver {
		t.Errorf("Expected ObjSaver to be %v, got %v", objectSaver, keeper.ObjSaver)
	}
	if keeper.objProvider != objectProvider {
		t.Errorf("Expected objProvider to be %v, got %v", objectProvider, keeper.objProvider)
	}
}

// Handles nil Logger gracefully
func TestNewKeeperService_NilLogger(t *testing.T) {
	provider := mocks.NewItemProvider(t)
	saver := mocks.NewItemSaver(t)
	objectSaver := mocks.NewObjectSaver(t)
	objectProvider := mocks.NewObjectProvider(t)

	keeper := NewKeeperService(nil, provider, saver, objectSaver, objectProvider)

	if keeper == nil {
		t.Fatal("Expected keeper to be initialized, got nil")
	}
	if keeper.log != nil {
		t.Errorf("Expected log to be nil, got %v", keeper.log)
	}
}

// Successfully creates a new item when valid data is provided
func TestCreateItem_Success(t *testing.T) {
	ctx := context.Background()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	item := &models.Item{
		Name:    "testItem",
		Content: []byte("testContent"),
		Version: uuid.New(),
		OwnerID: 1,
	}

	mockItemSaver := new(mocks.ItemSaver)
	mockItemSaver.On("Create", ctx, item).Return(item, nil)

	keeper := NewKeeperService(log, nil, mockItemSaver, nil, nil)

	createdItem, err := keeper.CreateItem(ctx, item)
	assert.NoError(t, err)
	assert.Equal(t, item, createdItem)
}

// Handles and logs errors when item creation fails
func TestCreateItem_Failure(t *testing.T) {
	ctx := context.Background()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	item := &models.Item{
		Name:    "testItem",
		Content: []byte("testContent"),
		Version: uuid.New(),
		OwnerID: 1,
	}

	mockItemSaver := new(mocks.ItemSaver)
	expectedErr := errors.New("creation failed")
	mockItemSaver.On("Create", ctx, item).Return(nil, expectedErr)

	keeper := NewKeeperService(log, nil, mockItemSaver, nil, nil)

	createdItem, err := keeper.CreateItem(ctx, item)
	assert.Error(t, err)
	assert.Nil(t, createdItem)
	assert.Equal(t, expectedErr, err)
}

// Successfully creates a file and returns its version string
func TestCreateFileSuccess(t *testing.T) {
	ctx := context.Background()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	itemProvider := mocks.NewItemProvider(t)
	itemSaver := mocks.NewItemSaver(t)
	objectSaver := mocks.NewObjectSaver(t)
	objectProvider := mocks.NewObjectProvider(t)

	keeper := NewKeeperService(log, itemProvider, itemSaver, objectSaver, objectProvider)

	file := &models.File{
		Item: models.Item{
			Name:    "testfile",
			Content: []byte("test content"),
			OwnerID: 1,
		},
		Size: 1024,
	}

	objectSaver.On("Create", ctx, file).Return("fileID123", nil)
	itemSaver.On("Create", ctx, &file.Item).Return(&models.Item{
		Name:    "testfile",
		Content: []byte("test content"),
		Version: uuid.New(),
		OwnerID: 1,
		FileID:  "fileID123",
	}, nil)

	version, err := keeper.CreateFile(ctx, file)

	assert.NoError(t, err)
	assert.NotEmpty(t, version)
}

// Handles error when ObjSaver.Create fails
func TestCreateFileObjSaverError(t *testing.T) {
	ctx := context.Background()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	itemProvider := mocks.NewItemProvider(t)
	itemSaver := mocks.NewItemSaver(t)
	objectSaver := mocks.NewObjectSaver(t)
	objectProvider := mocks.NewObjectProvider(t)

	keeper := NewKeeperService(log, itemProvider, itemSaver, objectSaver, objectProvider)

	file := &models.File{
		Item: models.Item{
			Name:    "testfile",
			Content: []byte("test content"),
			OwnerID: 1,
		},
		Size: 1024,
	}

	objectSaver.On("Create", ctx, file).Return("", fmt.Errorf("failed to create file"))

	version, err := keeper.CreateFile(ctx, file)

	assert.Error(t, err)
	assert.Empty(t, version)
}

// Successfully updates an item when valid data is provided
func TestUpdateItem_Success(t *testing.T) {
	ctx := context.Background()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	itemProvider := mocks.NewItemProvider(t)
	itemSaver := mocks.NewItemSaver(t)
	objectSaver := mocks.NewObjectSaver(t)
	objectProvider := mocks.NewObjectProvider(t)

	keeper := NewKeeperService(log, itemProvider, itemSaver, objectSaver, objectProvider)

	item := &models.Item{
		Name:    "testItem",
		Content: []byte("testContent"),
		Version: uuid.New(),
		OwnerID: 1,
		FileID:  "",
	}

	itemSaver.On("Update", ctx, item).Return(item, nil)

	updatedItem, err := keeper.UpdateItem(ctx, item)

	assert.NoError(t, err)
	assert.Equal(t, item, updatedItem)
	itemSaver.AssertExpectations(t)
}

// Successfully deletes an item when it exists
func TestDeleteItem_Success(t *testing.T) {
	ctx := context.Background()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	itemProvider := mocks.NewItemProvider(t)
	itemSaver := mocks.NewItemSaver(t)
	objectSaver := mocks.NewObjectSaver(t)
	objectProvider := mocks.NewObjectProvider(t)

	keeper := NewKeeperService(log, itemProvider, itemSaver, objectSaver, objectProvider)

	item := &models.Item{
		Name:    "testItem",
		OwnerID: 1,
		FileID:  "file123",
	}

	itemProvider.On("Get", ctx, "testItem", int64(1)).Return(item, nil)
	itemSaver.On("Delete", ctx, item).Return(nil)
	objectSaver.On("Delete", ctx, "file123").Return(nil)

	err := keeper.DeleteItem(ctx, item)

	assert.NoError(t, err)
	itemProvider.AssertExpectations(t)
	itemSaver.AssertExpectations(t)
	objectSaver.AssertExpectations(t)
}

// Item does not exist in the database
func TestDeleteItem_ItemNotFound(t *testing.T) {
	ctx := context.Background()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	itemProvider := mocks.NewItemProvider(t)
	itemSaver := mocks.NewItemSaver(t)
	objectSaver := mocks.NewObjectSaver(t)
	objectProvider := mocks.NewObjectProvider(t)

	keeper := NewKeeperService(log, itemProvider, itemSaver, objectSaver, objectProvider)

	item := &models.Item{
		Name:    "nonExistentItem",
		OwnerID: 1,
	}

	itemProvider.On("Get", ctx, "nonExistentItem", int64(1)).Return(nil, fmt.Errorf("item not found"))

	err := keeper.DeleteItem(ctx, item)

	assert.Error(t, err)
	assert.EqualError(t, err, "item not found")
	itemProvider.AssertExpectations(t)
}

// Handles the case where the item does not exist
func TestGetItem_NotFound(t *testing.T) {
	ctx := context.Background()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	itemProvider := mocks.NewItemProvider(t)
	itemSaver := mocks.NewItemSaver(t)
	objectSaver := mocks.NewObjectSaver(t)
	objectProvider := mocks.NewObjectProvider(t)

	keeper := NewKeeperService(log, itemProvider, itemSaver, objectSaver, objectProvider)

	itemProvider.On("Get", ctx, "nonExistentItem", int64(1)).Return(nil, fmt.Errorf("item not found"))

	item, err := keeper.GetItem(ctx, "nonExistentItem", int64(1))

	assert.Error(t, err)
	assert.Nil(t, item)
	assert.EqualError(t, err, "item not found")
	itemProvider.AssertExpectations(t)
}

// Successfully retrieves an item and its associated file
func TestGetFileSuccess(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)
	itemName := "testItem"
	fileID := "file123"
	version := uuid.New()

	itemProvider := mocks.NewItemProvider(t)
	objectProvider := mocks.NewObjectProvider(t)

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	expectedItem := &models.Item{
		Name:    itemName,
		Content: []byte("content"),
		Version: version,
		OwnerID: userID,
		FileID:  fileID,
	}

	expectedFile := &models.File{
		Item: *expectedItem,
		Size: 1024,
		Data: nil,
	}

	itemProvider.On("Get", ctx, itemName, userID).Return(expectedItem, nil)
	objectProvider.On("Get", ctx, fileID).Return(expectedFile, nil)

	keeper := NewKeeperService(log, itemProvider, nil, nil, objectProvider)

	file, err := keeper.GetFile(ctx, itemName, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedFile, file)

	itemProvider.AssertExpectations(t)
	objectProvider.AssertExpectations(t)
}

// Item retrieval fails due to non-existent item
func TestGetFileItemNotFound(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)
	itemName := "nonExistentItem"

	itemProvider := mocks.NewItemProvider(t)
	objectProvider := mocks.NewObjectProvider(t)

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	itemProvider.On("Get", ctx, itemName, userID).Return(nil, fmt.Errorf("item not found"))

	keeper := NewKeeperService(log, itemProvider, nil, nil, objectProvider)

	file, err := keeper.GetFile(ctx, itemName, userID)
	assert.Error(t, err)
	assert.Nil(t, file)
	assert.Contains(t, err.Error(), "failed to get info item for file")

	itemProvider.AssertExpectations(t)
}

// List items successfully when userID is valid and items exist
func TestListItemsSuccess(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)
	expectedItems := []*models.Item{
		{Name: "Item1", OwnerID: userID},
		{Name: "Item2", OwnerID: userID},
	}

	mockItemProvider := new(mocks.ItemProvider)
	mockItemProvider.On("List", ctx, userID).Return(expectedItems, nil)

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	keeper := NewKeeperService(log, mockItemProvider, nil, nil, nil)

	items, err := keeper.ListItems(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
	mockItemProvider.AssertExpectations(t)
}

// Handle error when itmProvider.List returns an error
func TestListItemsError(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)
	expectedError := errors.New("failed to list items")

	mockItemProvider := new(mocks.ItemProvider)
	mockItemProvider.On("List", ctx, userID).Return(nil, expectedError)

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	keeper := NewKeeperService(log, mockItemProvider, nil, nil, nil)

	items, err := keeper.ListItems(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Equal(t, expectedError, err)
	mockItemProvider.AssertExpectations(t)
}
