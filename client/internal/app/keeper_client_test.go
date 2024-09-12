package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	commonmocks "github.com/ajugalushkin/goph-keeper/mocks"
)

// Initializes KeeperClient with valid grpc.ClientConn
func TestInitializesKeeperClientWithValidClientConn(t *testing.T) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
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

	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
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

// Successfully opens the file specified by filePath
func TestCreateItemStream_Success(t *testing.T) {
	ctx := context.Background()
	name := "testfile.txt"
	content := []byte("test content")
	file, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	mockClient := commonmocks.NewKeeperServiceV1Client(t)
	mockStream := mocks.NewCreateItemStreamClient(t)
	mockClient.On("CreateItemStreamV1", ctx).Return(mockStream, nil)
	mockStream.On("Send", mock.Anything).Return(nil)
	mockStream.On("CloseAndRecv").Return(&keeperv1.CreateItemResponseV1{Name: name, Version: "1"}, nil)

	k := &KeeperClient{api: mockClient}
	res, err := k.CreateItemStream(ctx, name, file, content)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, name, res.GetName())
}

// File specified by filePath does not exist
func TestCreateItemStream_FileNotExist(t *testing.T) {
	ctx := context.Background()
	name := "nonexistentfile.txt"
	content := []byte("test content")

	file, err := os.Open(name)
	if err == nil {
		t.Fatalf("expected error when opening nonexistent file, got none")
	}

	mockClient := commonmocks.NewKeeperServiceV1Client(t)

	mockClient.On("CreateItemStreamV1", ctx).Return(nil, errors.New(""))

	k := &KeeperClient{api: mockClient}
	res, err := k.CreateItemStream(ctx, name, file, content)

	assert.Error(t, err)
	assert.Nil(t, res)
}

// Error occurs while sending file info to the server
func TestCreateItemStream_ErrorSendingFileInfo(t *testing.T) {
	ctx := context.Background()
	name := "testfile.txt"
	content := []byte("test content")
	file, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	mockClient := commonmocks.NewKeeperServiceV1Client(t)
	mockStream := mocks.NewCreateItemStreamClient(t)
	mockClient.On("CreateItemStreamV1", context.Background()).Return(mockStream, nil)
	mockStream.On("Send", mock.Anything).Return(errors.New("error sending file info"))
	mockStream.On("RecvMsg", nil).Return(errors.New("error receiving message"))

	k := &KeeperClient{api: mockClient}
	_, err = k.CreateItemStream(ctx, name, file, content)

	assert.Error(t, err)
}

// Error occurs while creating the stream
func TestCreateItemStream_ErrorOnStreamCreation(t *testing.T) {
	ctx := context.Background()
	name := "testfile.txt"
	content := []byte("test content")
	file, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	mockClient := commonmocks.NewKeeperServiceV1Client(t)
	mockClient.On("CreateItemStreamV1", ctx).Return(nil, errors.New("stream creation error"))

	k := &KeeperClient{api: mockClient}
	_, err = k.CreateItemStream(ctx, name, file, content)

	assert.Error(t, err)
}

// Error occurs while sending a chunk to the server
func TestCreateItemStream_ErrorSendingChunk(t *testing.T) {
	ctx := context.Background()
	name := "testfile.txt"
	content := []byte("test content")
	file, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	mockClient := commonmocks.NewKeeperServiceV1Client(t)
	mockStream := mocks.NewCreateItemStreamClient(t)
	mockClient.On("CreateItemStreamV1", context.Background()).Return(mockStream, nil)
	mockStream.On("Send", mock.Anything).Return(errors.New("error sending chunk"))
	mockStream.On("RecvMsg", nil).Return(errors.New("error receiving message"))

	k := &KeeperClient{api: mockClient}
	_, err = k.CreateItemStream(ctx, name, file, content)

	assert.Error(t, err)
}

// Error occurs while receiving the response after closing the stream
func TestCreateItemStream_ErrorOnResponseReceive(t *testing.T) {
	ctx := context.Background()
	name := "testfile.txt"
	content := []byte("test content")
	file, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	mockClient := commonmocks.NewKeeperServiceV1Client(t)
	mockStream := mocks.NewCreateItemStreamClient(t)
	mockClient.On("CreateItemStreamV1", ctx).Return(mockStream, nil)
	mockStream.On("Send", mock.Anything).Return(nil)
	mockStream.On("CloseAndRecv").Return(nil, errors.New("error receiving response"))

	k := &KeeperClient{api: mockClient}
	res, err := k.CreateItemStream(ctx, name, file, content)

	assert.Error(t, err)
	assert.Nil(t, res)
}

// Successfully updates an item when valid request is provided
func TestUpdateItem_Success(t *testing.T) {
	ctx := context.Background()
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
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

// Returns an error when the context is canceled or expired
func TestUpdateItem_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	req := &keeperv1.UpdateItemRequestV1{
		Name:    "test-item",
		Content: []byte("test-content"),
	}

	mockAPI.On("UpdateItemV1", ctx, req).Return(nil, context.Canceled)

	result, err := client.UpdateItem(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, context.Canceled, errors.Unwrap(err))

	mockAPI.AssertExpectations(t)
}

// Handles grpc errors and returns a formatted error message
func TestUpdateItem_GrpcError(t *testing.T) {
	ctx := context.Background()
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	req := &keeperv1.UpdateItemRequestV1{
		Name:    "test-item",
		Content: []byte("test-content"),
	}
	expectedErr := errors.New("grpc error")

	mockAPI.On("UpdateItemV1", ctx, req).Return(nil, expectedErr)

	_, err := client.UpdateItem(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "keeper.UpdateItem: grpc error")

	mockAPI.AssertExpectations(t)
}

// Handles network failures gracefully
func TestUpdateItem_NetworkFailure(t *testing.T) {
	ctx := context.Background()
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
	client := &KeeperClient{api: mockAPI}

	req := &keeperv1.UpdateItemRequestV1{
		Name:    "test-item",
		Content: []byte("test-content"),
	}

	mockAPI.On("UpdateItemV1", ctx, req).Return(nil, errors.New("network error"))

	result, err := client.UpdateItem(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, result)

	mockAPI.AssertExpectations(t)
}

// Successfully deletes an item when valid request is provided
func TestDeleteItem_Success(t *testing.T) {
	ctx := context.Background()
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
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
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
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
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
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
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
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

// Successfully retrieves a stream when given a valid context and name
func TestGetFile_Success(t *testing.T) {
	ctx := context.Background()
	name := "test-file"

	mockStream := commonmocks.NewServerStreamingClient[keeperv1.GetItemStreamResponseV1](t)
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
	mockAPI.On("GetItemStreamV1", ctx, &keeperv1.GetItemRequestV1{Name: name}).Return(mockStream, nil)

	client := &KeeperClient{api: mockAPI}

	stream, err := client.GetFile(ctx, name)

	assert.NoError(t, err)
	assert.Equal(t, mockStream, stream)
}

// Handles the case where the context is canceled or expired
func TestGetFile_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	name := "testfile"

	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
	mockAPI.On("GetItemStreamV1", ctx, &keeperv1.GetItemRequestV1{Name: name}).Return(nil, context.Canceled)

	client := &KeeperClient{api: mockAPI}

	stream, err := client.GetFile(ctx, name)

	assert.Error(t, err)
	assert.Nil(t, stream)
	assert.Equal(t, context.Canceled, errors.Unwrap(err))
}

// Handles network failures or interruptions during the API call
func TestGetFile_NetworkFailure(t *testing.T) {
	ctx := context.Background()
	name := "testfile"

	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
	mockAPI.On("GetItemStreamV1", ctx, &keeperv1.GetItemRequestV1{Name: name}).Return(nil, errors.New("network error"))

	client := &KeeperClient{api: mockAPI}

	stream, err := client.GetFile(ctx, name)

	assert.Error(t, err)
	assert.Nil(t, stream)
}

// Properly formats and returns errors when the API call fails
func TestGetFile_Error(t *testing.T) {
	ctx := context.Background()
	name := "testfile"

	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
	mockAPI.On("GetItemStreamV1", ctx, &keeperv1.GetItemRequestV1{Name: name}).Return(nil, errors.New("API call failed"))

	client := &KeeperClient{api: mockAPI}

	_, err := client.GetFile(ctx, name)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client.keeper.GetItem")
}

// Successfully retrieves a list of items when the API call returns a valid response
func TestListItems_Success(t *testing.T) {
	ctx := context.Background()
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
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
	mockAPI := commonmocks.NewKeeperServiceV1Client(t)
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

// Returns a valid grpc.ClientConn when provided a valid token
func TestGetKeeperConnectionWithValidToken(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	address := "localhost:50051"
	token := "valid-token"

	conn := GetKeeperConnection(log, address, token)

	if conn == nil {
		t.Fatalf("Expected a valid grpc.ClientConn, got nil")
	}
}

// Handles nil token input gracefully
func TestGetKeeperConnectionWithNilToken(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	address := "invalid_connection"
	token := ""

	conn := GetKeeperConnection(log, address, token)

	if conn.Target() != "invalid_connection" {
		t.Fatalf("Expected nil grpc.ClientConn, got a valid connection")
	}
}

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
