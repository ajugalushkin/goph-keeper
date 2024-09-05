package v1

import (
	"context"
	"errors"

	"github.com/ajugalushkin/goph-keeper/server/internal/services"

	//"io"
	//"log"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage"
	//"github.com/ajugalushkin/goph-keeper/server/internal/storage/minio"
)

type Keeper interface {
	CreateItem(
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
	ListItems(
		ctx context.Context,
		userID int64,
	) (list []*models.Item, err error)
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

func (s *serverAPI) CreateStream(
	stream keeperv1.KeeperServiceV1_CreateItemStreamV1Server,
) error {
	//req, err := stream.Recv()
	//if err != nil {
	//	//	return logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	//}
	////
	//fileName := req.GetInfo().GetName()
	//fileType := req.GetInfo().GetType()
	////log.Printf("receive an upload-image request for laptop %s with image type %s", laptopID, imageType)
	////
	////laptop, err := server.laptopStore.Find(laptopID)
	////if err != nil {
	////	return logError(status.Errorf(codes.Internal, "cannot find laptop: %v", err))
	////}
	////if laptop == nil {
	////	return logError(status.Errorf(codes.InvalidArgument, "laptop id %s doesn't exist", laptopID))
	////}
	////
	////imageData := bytes.Buffer{}
	////imageSize := 0
	//
	//storage, err := minio.NewMinioClient()
	//if err != nil {
	//	return err
	//}
	//
	//for {
	//	//err := contextError(stream.Context())
	//	//if err != nil {
	//	//	return err
	//	//}
	//	//
	//	//log.Print("waiting to receive more data")
	//
	//	req, err := stream.Recv()
	//	if err != nil {
	//		if err == io.EOF {
	//			//log.Print("no more data")
	//			break
	//		}
	//		//return logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
	//	}
	//
	//	chunk := req.GetChunkData()
	//	size := len(chunk)
	//
	//	log.Printf("received a chunk with size: %d", size)
	//
	//	//imageSize += size
	//	//if imageSize > maxImageSize {
	//	//	return logError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize))
	//	//}
	//
	//	////write slowly
	//	////time.Sleep(time.Second)
	//
	//	//_, err = imageData.Write(chunk)
	//	//if err != nil {
	//	//	return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
	//	//}
	//}

	//imageID, err := server.imageStore.Save(laptopID, imageType, imageData)
	//if err != nil {
	//	return logError(status.Errorf(codes.Internal, "cannot save image to the store: %v", err))
	//}
	//
	//res := &pb.UploadImageResponse{
	//	Id:   imageID,
	//	Size: uint32(imageSize),
	//}
	//
	//err = stream.SendAndClose(res)
	//if err != nil {
	//	return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	//}
	//
	//log.Printf("saved image with id: %s, size: %d", imageID, imageSize)
	//return nil
	return nil
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

func (s *serverAPI) DeleteItem(
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
