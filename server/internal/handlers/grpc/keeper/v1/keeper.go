package v1

import (
	"context"
	"strconv"

	"google.golang.org/grpc"

	"github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

type Keeper interface {
	ListItem(
		ctx context.Context,
		since int64,
	) (list *models.ListItem, err error)
	SaveItem(ctx context.Context, item *models.Item) (serverUpdateAt string, err error)
}

type serverAPI struct {
	v1.UnimplementedKeeperServiceV1Server
	keeper Keeper
}

func Register(gRPC *grpc.Server, keeper Keeper) {
	v1.RegisterKeeperServiceV1Server(gRPC, &serverAPI{
		keeper: keeper,
	})
}

func (s *serverAPI) ListItemsV1(
	ctx context.Context,
	req *v1.ListItemRequestV1,
) (*v1.ListItemResponseV1, error) {
	_, err := s.keeper.ListItem(ctx, req.GetSince())
	if err != nil {
		return nil, err
	}

	return &v1.ListItemResponseV1{}, nil
}

func (s *serverAPI) SetItemV1(
	ctx context.Context,
	req *v1.SetItemRequestV1,
) (*v1.SetItemResponseV1, error) {
	_, err := s.keeper.SaveItem(ctx, &models.Item{
		ID:              req.GetItem().GetId(),
		Name:            req.GetItem().GetName(),
		Type:            req.GetItem().GetType().String(),
		Value:           req.GetItem().GetValue(),
		ServerUpdatedAt: strconv.FormatInt(req.GetItem().ServerUpdatedAt, 10),
		IsDeleted:       req.GetItem().IsDeleted,
	})
	if err != nil {
		return nil, err
	}

	return &v1.SetItemResponseV1{
		//ServerUpdatedAt: updatedAt,
	}, nil
}
