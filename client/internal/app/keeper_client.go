package app

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

type KeeperClient struct {
	api keeperv1.KeeperServiceV1Client
}

// NewKeeperClient returns a new keeper client
func NewKeeperClient(cc *grpc.ClientConn) *KeeperClient {
	service := keeperv1.NewKeeperServiceV1Client(cc)
	return &KeeperClient{service}
}

func (k *KeeperClient) ListItem(ctx context.Context, since int64) (error, *keeperv1.ListItemResponseV1) {
	const op = "client.keeper.Register"

	list, err := k.api.ListItemV1(ctx, &keeperv1.ListItemRequestV1{
		Since: since,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err), nil
	}

	return nil, list
}

func (k *KeeperClient) SetItem(ctx context.Context, item *keeperv1.Item) (int64, error) {
	const op = "client.keeper.Login"

	resp, err := k.api.SetItemV1(ctx, &keeperv1.SetItemRequestV1{
		Item: item,
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetServerUpdatedAt(), nil
}
