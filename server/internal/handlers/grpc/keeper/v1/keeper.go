package v1

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	"github.com/ajugalushkin/goph-keeper/server/internal/services/interceptors"
)

type Keeper interface {
	CreateItem(
		ctx context.Context,
		item *models.Item,
	) (*models.Item, error)
	ListItem(
		ctx context.Context,
		since int64,
	) (list *models.ListItem, err error)
}

type serverAPI struct {
	keeperv1.UnimplementedKeeperServiceV1Server
	keeper Keeper
}

func Register(gRPC *grpc.Server, keeper Keeper) {
	keeperv1.RegisterKeeperServiceV1Server(gRPC, &serverAPI{
		keeper: keeper,
	})
}

func (s *serverAPI) CreateSecret(
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

	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
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
		//if errors.Is(err, storage.ErrSecretConflict) {
		//	return nil, status.Error(codes.AlreadyExists, "secret already exists")
		//}
		return nil, status.Error(codes.Internal, "failed to create secret")
	}

	return &keeperv1.CreateItemResponseV1{
		Name:    req.GetName(),
		Version: item.Version.String(),
	}, nil
}

func (s *serverAPI) ListItemsV1(
	ctx context.Context,
	req *keeperv1.ListItemRequestV1,
) (*keeperv1.ListItemResponseV1, error) {
	_, err := s.keeper.ListItem(ctx, req.GetSince())
	if err != nil {
		return nil, err
	}

	return &keeperv1.ListItemResponseV1{}, nil
}

func (s *serverAPI) SetItemV1(
	ctx context.Context,
	req *keeperv1.SetItemRequestV1,
) (*keeperv1.SetItemResponseV1, error) {
	//_, err := s.keeper.SaveItem(ctx, &models.Item{
	//	ID:              req.GetItem().GetId(),
	//	Name:            req.GetItem().GetName(),
	//	Type:            req.GetItem().GetType().String(),
	//	Value:           req.GetItem().GetValue(),
	//	ServerUpdatedAt: strconv.FormatInt(req.GetItem().ServerUpdatedAt, 10),
	//	IsDeleted:       req.GetItem().IsDeleted,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &keeperv1.SetItemResponseV1{
	//	//ServerUpdatedAt: updatedAt,
	//}, nil
	return nil, nil
}
