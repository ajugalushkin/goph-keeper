package v1

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ajugalushkin/goph-keeper/server/interceptors"

	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage"
)

//go:generate mockery --name Keeper
type Keeper interface {
	CreateItem(
		ctx context.Context,
		item *models.Item,
	) (*models.Item, error)
	CreateFile(
		ctx context.Context,
		file *models.File,
	) (string, error)
	UpdateItem(
		ctx context.Context,
		item *models.Item,
	) (*models.Item, error)
	DeleteItem(
		ctx context.Context,
		item *models.Item,
	) error
	GetItem(
		ctx context.Context,
		name string,
		userID int64,
	) (*models.Item, error)
	GetFile(
		ctx context.Context,
		name string,
		userID int64,
	) (*models.File, error)
	ListItems(
		ctx context.Context,
		userID int64,
	) (list []*models.Item, err error)
}

//go:generate mockery --name CreateItemStreamServer
type CreateItemStreamServer interface {
	keeperv1.KeeperServiceV1_CreateItemStreamV1Server
}

type serverAPI struct {
	keeperv1.UnimplementedKeeperServiceV1Server
	keeper Keeper
}

// Register registers the KeeperServiceV1 server to the gRPC server.
// It takes a gRPC server instance and a Keeper interface as parameters.
// The Keeper interface provides methods for interacting with the storage layer.
// The serverAPI struct implements the KeeperServiceV1Server interface,
// which is responsible for handling the gRPC requests.
func Register(
	gRPC *grpc.Server,
	keeper Keeper,
) {
	keeperv1.RegisterKeeperServiceV1Server(gRPC, &serverAPI{
		keeper: keeper,
	})
}

// CreateItemV1 handles the creation of a new item in the storage.
// It validates the input request, retrieves the user ID from the context,
// creates a new item model, and calls the Keeper's CreateItem method.
// If the item already exists, it returns an AlreadyExists error.
// If any other error occurs, it returns an Internal error.
// Otherwise, it returns the created item's name and version.
func (s *serverAPI) CreateItemV1(
	ctx context.Context,
	req *keeperv1.CreateItemRequestV1,
) (*keeperv1.CreateItemResponseV1, error) {
	// Validate the input request.
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}
	if len(req.GetContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty secret content")
	}

	// Retrieve the user ID from the context.
	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	// Create a new item model.
	item, err := s.keeper.CreateItem(ctx, &models.Item{
		Name:    req.GetName(),
		Content: req.GetContent(),
		Version: uuid.UUID{},
		OwnerID: userID,
	})
	if err != nil {
		// Handle specific errors.
		if errors.Is(err, storage.ErrItemConflict) {
			return nil, status.Error(codes.AlreadyExists, "item already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create item")
	}

	// Return the created item's name and version.
	return &keeperv1.CreateItemResponseV1{
		Name:    req.GetName(),
		Version: item.Version.String(),
	}, nil
}

// CreateItemStreamV1 handles the streaming of file data for creating a new item in the storage.
// It receives the file data in chunks and writes it to a temporary file.
// Once all the data is received, it creates a new file model and calls the Keeper's CreateFile method.
// If any error occurs during the process, it returns an appropriate error.
//
// Parameters:
// - stream: A KeeperServiceV1_CreateItemStreamV1Server instance that handles the streaming of file data.
//
// Return Value:
// - error: An error if any error occurs during the process. Otherwise, it returns nil.
func (s *serverAPI) CreateItemStreamV1(
	stream keeperv1.KeeperServiceV1_CreateItemStreamV1Server,
) error {
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive file info"))
	}

	nameItem := req.GetInfo().GetName()
	fileContent := req.GetInfo().GetContent()
	log.Printf("receive an upload-file request for %s", nameItem)

	ctx := stream.Context()
	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int64)
	if !ok || userID == 0 {
		return logError(status.Errorf(codes.Unauthenticated, "empty user id"))
	}

	fileSize := 0

	tempFile, err := os.CreateTemp("", nameItem)
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Print("no more data")
				break
			}
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("received a chunk with size: %d", size)

		fileSize += size

		_, err = tempFile.Write(chunk)
		if err != nil {
			return err
		}
	}

	version, err := s.keeper.CreateFile(context.Background(),
		&models.File{
			Item: models.Item{
				Name:    nameItem,
				Content: fileContent,
				Version: uuid.UUID{},
				OwnerID: userID,
			},
			Size: int64(fileSize),
			Data: tempFile,
		})
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot write file data: %v", err))
	}

	res := &keeperv1.CreateItemResponseV1{
		Name:    nameItem,
		Version: version,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	log.Printf("saved file with id: %s, file size: %d, version: %s", nameItem, fileSize, version)

	return nil
}

// UpdateItemV1 handles the update of an existing item in the storage.
// It validates the input request, retrieves the user ID from the context,
// updates the item model, and calls the Keeper's UpdateItem method.
// If the item does not exist, it returns a NotFound error.
// If any other error occurs, it returns an Internal error.
// Otherwise, it returns the updated item's name and version.
func (s *serverAPI) UpdateItemV1(
	ctx context.Context,
	req *keeperv1.UpdateItemRequestV1,
) (*keeperv1.UpdateItemResponseV1, error) {
	// Validate the input request.
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}
	if len(req.GetContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty secret content")
	}

	// Retrieve the user ID from the context.
	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	// Update the item model.
	item, err := s.keeper.UpdateItem(ctx, &models.Item{
		Name:    req.GetName(),
		Content: req.GetContent(),
		OwnerID: userID,
	})
	if err != nil {
		// Handle specific errors.
		if errors.Is(err, storage.ErrItemNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}
		return nil, status.Error(codes.Internal, "failed to create secret")
	}
	// Return the updated item's name and version.
	return &keeperv1.UpdateItemResponseV1{
		Name:    req.GetName(),
		Version: item.Version.String(),
	}, nil
}

// DeleteItemV1 handles the deletion of an existing item in the storage.
// It validates the input request, retrieves the user ID from the context,
// creates a new item model, and calls the Keeper's DeleteItem method.
// If the item does not exist, it returns a NotFound error.
// If any other error occurs, it returns an Internal error.
// Otherwise, it returns the deleted item's name.
func (s *serverAPI) DeleteItemV1(
	ctx context.Context,
	request *keeperv1.DeleteItemRequestV1,
) (*keeperv1.DeleteItemResponseV1, error) {
	// Validate the input request.
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}

	// Retrieve the user ID from the context.
	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	// Create a new item model.
	item := &models.Item{
		Name:    request.GetName(),
		OwnerID: userID,
	}

	// Call the Keeper's DeleteItem method.
	err := s.keeper.DeleteItem(ctx, item)
	if err != nil {
		// Handle specific errors.
		if errors.Is(err, storage.ErrItemNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}
		return nil, status.Error(codes.Internal, "failed to del secret")
	}

	// Return the deleted item's name.
	return &keeperv1.DeleteItemResponseV1{
		Name: request.GetName(),
	}, nil
}

// GetItemV1 retrieves an item from the storage based on the provided name and user ID.
//
// Parameters:
// - ctx: A context that carries deadlines, cancellation signals, and other request-scoped values.
// - request: A GetItemRequestV1 object containing the name of the item to retrieve.
//
// Return Value:
// - A GetItemResponseV1 object containing the retrieved item's name, content, and version.
// - An error if any error occurs during the retrieval process.
//
// If the provided name is empty, it returns an InvalidArgument error with a message indicating that the secret name is empty.
// If the user ID is not found in the context, it returns an Unauthenticated error with a message indicating that the empty user id is provided.
// If the item is not found in the storage, it returns a NotFound error with a message indicating that the item not found.
// If any other error occurs during the retrieval process, it returns an Internal error with a message indicating that the failed to get item.
func (s *serverAPI) GetItemV1(
	ctx context.Context,
	request *keeperv1.GetItemRequestV1,
) (*keeperv1.GetItemResponseV1, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret name is empty")
	}

	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	item, err := s.keeper.GetItem(ctx, request.GetName(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "item not found")
		}
		return nil, status.Error(codes.Internal, "failed to get item")
	}

	return &keeperv1.GetItemResponseV1{
		Name:    item.Name,
		Content: item.Content,
		Version: item.Version.String(),
	}, nil
}

// GetItemStreamV1 handles the streaming of file data for retrieving an item from the storage.
// It receives the file data in chunks and sends it to the client in a streaming manner.
//
// Parameters:
// - req: A GetItemRequestV1 object containing the name of the item to retrieve.
// - stream: A KeeperServiceV1_GetItemStreamV1Server instance that handles the streaming of file data.
//
// Return Value:
// - error: An error if any error occurs during the process. Otherwise, it returns nil.
//
// The function performs the following steps:
// 1. Validates the input request and retrieves the user ID from the context.
// 2. Calls the Keeper's GetFile method to retrieve the file data based on the provided name and user ID.
// 3. Sends the file content as the first response using the stream.
// 4. Reads and sends the file data in chunks until all the data is sent.
// 5. Returns any error that occurs during the process.
func (s *serverAPI) GetItemStreamV1(
	req *keeperv1.GetItemRequestV1,
	stream keeperv1.KeeperServiceV1_GetItemStreamV1Server,
) error {
	const op = "keeperv1.GetItemStreamV1"

	name := req.GetName()

	ctx := stream.Context()
	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int64)
	if !ok || userID == 0 {
		return logError(status.Errorf(codes.Unauthenticated, "op: %s, empty user id", op))
	}

	file, err := s.keeper.GetFile(context.Background(), name, userID)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "op: %s, file not found: %v", op, err))
	}
	defer file.Data.Close()

	err = stream.Send(&keeperv1.GetItemStreamResponseV1{
		Content: file.Item.Content,
	})
	if err != nil {
		return logError(status.Errorf(codes.Internal, "op: %s, cannot send info response: %v", op, err))
	}

	buff := make([]byte, 1024)
	for {
		bytesRead, err := file.Data.Read(buff)
		if err != nil {
			if err == io.EOF {
				log.Printf("return file: %s", name)
				break
			}
			return logError(status.Errorf(codes.Internal, "op: %s, cannot read file: %v", op, err))
		}

		err = stream.Send(&keeperv1.GetItemStreamResponseV1{
			ChunkData: buff[:bytesRead],
		})
		if err != nil {
			return logError(status.Errorf(codes.Internal, "op: %s, cannot send chunk data: %v", op, err))
		}
	}

	return nil
}

// ListItemsV1 retrieves a list of items from the storage based on the provided user ID.
//
// Parameters:
// - ctx: A context that carries deadlines, cancellation signals, and other request-scoped values.
// - req: A ListItemsRequestV1 object containing the user ID for which to retrieve items.
//
// Return Value:
// - A ListItemsResponseV1 object containing a list of SecretInfo objects representing the retrieved items.
// - An error if any error occurs during the retrieval process.
//
// If the user ID is not found in the context, it returns an Unauthenticated error with a message indicating that the empty user id is provided.
// If any error occurs while listing the items from the storage, it returns an Internal error with a message indicating that the failed to list secrets.
func (s *serverAPI) ListItemsV1(
	ctx context.Context,
	req *keeperv1.ListItemsRequestV1,
) (*keeperv1.ListItemsResponseV1, error) {
	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	secrets, err := s.keeper.ListItems(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list secrets")
	}

	keeperSecrets := make([]*keeperv1.SecretInfo, 0, len(secrets))
	for _, secret := range secrets {
		keeperSecrets = append(keeperSecrets, &keeperv1.SecretInfo{
			Name:    secret.Name,
			Content: secret.Content,
			Version: secret.Version.String(),
		})
	}
	return &keeperv1.ListItemsResponseV1{
		Secrets: keeperSecrets,
	}, nil
}

// contextError checks the context for any errors related to cancellation or deadline exceeded.
// It returns an error if the context is canceled or if the deadline is exceeded.
//
// Parameters:
// ctx: A context that carries deadlines, cancellation signals, and other request-scoped values.
//
// Return Value:
// An error if the context is canceled or if the deadline is exceeded.
// nil if there are no errors related to cancellation or deadline exceeded.
func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}

// logError is a helper function that logs an error if it is not nil.
// It returns the error unchanged.
//
// Parameters:
// - err: An error to be logged. If this parameter is nil, the function does nothing.
//
// Return Value:
// - error: The same error that was passed as a parameter.
func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
