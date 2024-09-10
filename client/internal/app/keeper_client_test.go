package app

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/gen/keeper/v1/mocks"
)

// Initializes KeeperClient with valid grpc.ClientConn
func TestInitializesKeeperClientWithValidClientConn(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := NewKeeperClient(conn)
	if client == nil {
		t.Fatalf("Expected non-nil KeeperClient")
	}
}

// Handles nil grpc.ClientConn gracefully
func TestHandlesNilClientConnGracefully(t *testing.T) {
	client := NewKeeperClient(nil)
	if client == nil {
		t.Fatalf("Expected non-nil KeeperClient")
	}
}

// Successfully creates an item when valid request is provided
func TestCreateItem_Success(t *testing.T) {
	ctx := context.Background()
	mockAPI := mocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	req := &keeperv1.CreateItemRequestV1{
		Name:    "test-item",
		Content: []byte("test-content"),
	}
	resp := &keeperv1.CreateItemResponseV1{
		Name:    "test-item",
		Version: "v1",
	}

	mockAPI.On("CreateItemV1", ctx, req).Return(resp, nil)

	result, err := client.CreateItem(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, resp, result)

	mockAPI.AssertExpectations(t)
}

// Returns an error when the context is canceled
func TestCreateItem_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mockAPI := mocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	req := &keeperv1.CreateItemRequestV1{
		Name:    "test-item",
		Content: []byte("test-content"),
	}

	mockAPI.On("CreateItemV1", ctx, req).Return(nil, context.Canceled)

	result, err := client.CreateItem(ctx, req)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, context.Canceled, err)

	mockAPI.AssertExpectations(t)
}

// stream

// File specified by filePath does not exist
func TestCreateItemStream_FileDoesNotExist(t *testing.T) {
	ctx := context.Background()
	name := "testfile"
	filePath := "nonexistentfile.txt"
	content := []byte("test content")

	mockClient := mocks.NewKeeperServiceV1Client(t)

	k := &KeeperClient{api: mockClient}

	res, err := k.CreateItemStream(ctx, name, filePath, content)
	assert.Error(t, err)
	assert.Nil(t, res)
}

// Successfully updates an item when valid request is provided
func TestUpdateItem_Success(t *testing.T) {
	ctx := context.Background()
	mockAPI := mocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	req := &keeperv1.UpdateItemRequestV1{
		Name:    "test-item",
		Content: []byte("test-content"),
	}
	resp := &keeperv1.UpdateItemResponseV1{
		Name:    "test-item",
		Version: "v1",
	}

	mockAPI.On("UpdateItemV1", ctx, req).Return(resp, nil)

	result, err := client.UpdateItem(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, resp, result)

	mockAPI.AssertExpectations(t)
}

// Successfully deletes an item when valid request is provided
func TestDeleteItem_Success(t *testing.T) {
	ctx := context.Background()
	mockAPI := mocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	req := &keeperv1.DeleteItemRequestV1{Name: "test-item"}
	resp := &keeperv1.DeleteItemResponseV1{Name: "test-item"}

	mockAPI.On("DeleteItemV1", ctx, req).Return(resp, nil)

	result, err := client.DeleteItem(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, resp, result)
	mockAPI.AssertExpectations(t)
}

// Handles network failures gracefully
func TestDeleteItem_NetworkFailure(t *testing.T) {
	ctx := context.Background()
	mockAPI := mocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	req := &keeperv1.DeleteItemRequestV1{Name: "test-item"}
	expectedErr := fmt.Errorf("network error")

	mockAPI.On("DeleteItemV1", ctx, req).Return(nil, expectedErr)

	result, err := client.DeleteItem(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "network error")
	mockAPI.AssertExpectations(t)
}

// Successfully retrieves item when valid request is provided
func TestGetItem_Success(t *testing.T) {
	ctx := context.Background()
	mockAPI := mocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	req := &keeperv1.GetItemRequestV1{Name: "test-item"}
	expectedResp := &keeperv1.GetItemResponseV1{Name: "test-item", Content: []byte("content"), Version: "v1"}

	mockAPI.On("GetItemV1", ctx, req).Return(expectedResp, nil)

	resp, err := client.GetItem(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	mockAPI.AssertExpectations(t)
}

// Handles error when API call fails
func TestGetItem_Error(t *testing.T) {
	ctx := context.Background()
	mockAPI := mocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	req := &keeperv1.GetItemRequestV1{Name: "test-item"}
	expectedErr := fmt.Errorf("API error")

	mockAPI.On("GetItemV1", ctx, req).Return(nil, expectedErr)

	resp, err := client.GetItem(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "client.keeper.GetItem")

	mockAPI.AssertExpectations(t)
}

// stream get file

// Successfully retrieves a list of items when the API call returns a valid response
func TestListItems_Success(t *testing.T) {
	ctx := context.Background()
	mockAPI := mocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	request := &keeperv1.ListItemsRequestV1{}
	expectedResponse := &keeperv1.ListItemsResponseV1{
		Secrets: []*keeperv1.SecretInfo{
			{Name: "Secret1"},
			{Name: "Secret2"},
		},
	}

	mockAPI.On("ListItemsV1", ctx, request).Return(expectedResponse, nil)

	response, err := client.ListItems(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
	mockAPI.AssertExpectations(t)
}

// Handles the case where the API call returns an error
func TestListItems_Error(t *testing.T) {
	ctx := context.Background()
	mockAPI := mocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	request := &keeperv1.ListItemsRequestV1{}
	expectedError := fmt.Errorf("some error")

	mockAPI.On("ListItemsV1", ctx, request).Return(nil, expectedError)

	response, err := client.ListItems(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, expectedError, errors.Unwrap(err))
	mockAPI.AssertExpectations(t)
}

// get connection

// Returns a map with correct method names as keys
func TestAuthMethodsReturnsCorrectKeys(t *testing.T) {
	expectedKeys := []string{
		keeperv1.KeeperServiceV1_ListItemsV1_FullMethodName,
		keeperv1.KeeperServiceV1_GetItemV1_FullMethodName,
		keeperv1.KeeperServiceV1_CreateItemV1_FullMethodName,
		keeperv1.KeeperServiceV1_CreateItemStreamV1_FullMethodName,
		keeperv1.KeeperServiceV1_DeleteItemV1_FullMethodName,
		keeperv1.KeeperServiceV1_UpdateItemV1_FullMethodName,
	}

	authMap := authMethods()

	for _, key := range expectedKeys {
		assert.True(t, authMap[key], "Expected key %s to be present in the map", key)
	}
}
