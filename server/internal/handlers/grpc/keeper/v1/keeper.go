package v1

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/ajugalushkin/goph-keeper/server/internal/services"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func Register(
	gRPC *grpc.Server,
	keeper Keeper,
) {
	keeperv1.RegisterKeeperServiceV1Server(gRPC, &serverAPI{
		keeper: keeper,
	})
}

func (s *serverAPI) CreateItemV1(
	ctx context.Context,
	req *keeperv1.CreateItemRequestV1,
) (*keeperv1.CreateItemResponseV1, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, ok := ctx.Value(services.ContextKeyUserID).(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	item, err := s.keeper.CreateItem(ctx, &models.Item{
		Name:    req.GetName(),
		Content: req.GetContent(),
		Version: uuid.UUID{},
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, storage.ErrItemConflict) {
			return nil, status.Error(codes.AlreadyExists, "item already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create item")
	}

	return &keeperv1.CreateItemResponseV1{
		Name:    req.GetName(),
		Version: item.Version.String(),
	}, nil
}

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
	userID, ok := ctx.Value(services.ContextKeyUserID).(int64)
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

func (s *serverAPI) UpdateItemV1(
	ctx context.Context,
	req *keeperv1.UpdateItemRequestV1,
) (*keeperv1.UpdateItemResponseV1, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}
	if len(req.GetContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty secret content")
	}

	userID, ok := ctx.Value(services.ContextKeyUserID).(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	item, err := s.keeper.UpdateItem(ctx, &models.Item{
		Name:    req.GetName(),
		Content: req.GetContent(),
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}
		return nil, status.Error(codes.Internal, "failed to create secret")
	}
	return &keeperv1.UpdateItemResponseV1{
		Name:    req.GetName(),
		Version: item.Version.String(),
	}, nil
}

func (s *serverAPI) DeleteItemV1(
	ctx context.Context,
	request *keeperv1.DeleteItemRequestV1,
) (*keeperv1.DeleteItemResponseV1, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}

	userID, ok := ctx.Value(services.ContextKeyUserID).(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	item := &models.Item{
		Name:    request.GetName(),
		OwnerID: userID,
	}

	err := s.keeper.DeleteItem(ctx, item)
	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}
		return nil, status.Error(codes.Internal, "failed to create secret")
	}

	return &keeperv1.DeleteItemResponseV1{
		Name: request.GetName(),
	}, nil
}

func (s *serverAPI) GetItemV1(
	ctx context.Context,
	request *keeperv1.GetItemRequestV1,
) (*keeperv1.GetItemResponseV1, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret name is empty")
	}

	userID, ok := ctx.Value(services.ContextKeyUserID).(int64)
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

func (s *serverAPI) GetItemStreamV1(
	req *keeperv1.GetItemRequestV1,
	stream keeperv1.KeeperServiceV1_GetItemStreamV1Server,
) error {
	const op = "keeperv1.GetItemStreamV1"

	name := req.GetName()

	ctx := stream.Context()
	userID, ok := ctx.Value(services.ContextKeyUserID).(int64)
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

func (s *serverAPI) ListItemsV1(
	ctx context.Context,
	req *keeperv1.ListItemsRequestV1,
) (*keeperv1.ListItemsResponseV1, error) {
	userID, ok := ctx.Value(services.ContextKeyUserID).(int64)
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

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
