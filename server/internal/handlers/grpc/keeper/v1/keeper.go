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
	//CreateItemStream(
	//	ctx context.Context,
	//	stream  )
	CreateItem(
		ctx context.Context,
		item *models.Item,
	) (*models.Item, error)
	GetItem(
		ctx context.Context,
		name string,
		userID int64,
	) (*models.Item, error)
	//ListItem(
	//	ctx context.Context,
	//	since int64,
	//) (list *models.ListItem, err error)
	//Save(ctx context.Context)
	//SaveStream(ctx context.Context)
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

//func (s *serverAPI) ListItemsV1(
//	ctx context.Context,
//	req *keeperv1.ListItemRequestV1,
//) (*keeperv1.ListItemResponseV1, error) {
//	_, err := s.keeper.ListItem(ctx, req.GetSince())
//	if err != nil {
//		return nil, err
//	}
//
//	return &keeperv1.ListItemResponseV1{}, nil
//}
//
//func (s *serverAPI) SetItemV1(
//	ctx context.Context,
//	req *keeperv1.SetItemRequestV1,
//) (*keeperv1.SetItemResponseV1, error) {
//	//_, err := s.keeper.SaveItem(ctx, &models.Item{
//	//	ID:              req.GetItem().GetId(),
//	//	Name:            req.GetItem().GetName(),
//	//	Type:            req.GetItem().GetType().String(),
//	//	Value:           req.GetItem().GetValue(),
//	//	ServerUpdatedAt: strconv.FormatInt(req.GetItem().ServerUpdatedAt, 10),
//	//	IsDeleted:       req.GetItem().IsDeleted,
//	//})
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//return &keeperv1.SetItemResponseV1{
//	//	//ServerUpdatedAt: updatedAt,
//	//}, nil
//	return nil, nil
//}
