package app

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	authv1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
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

//func (k *KeeperClient) ListItemV1(ctx context.Context, in *keeperv1.ListItemRequestV1, opts ...grpc.CallOption) (*keeperv1.ListItemResponseV1, error){
//	return k.service.ListItemV1(ctx, in, opts...)
//}
//SetItemV1(ctx context.Context, in *SetItemRequestV1, opts ...grpc.CallOption) (*SetItemResponseV1, error)

func (k *KeeperClient) ListItem(ctx context.Context, email string, password string) error {
	const op = "client.keeper.Register"

	_, err := k.api.ListItemV1(ctx, &keeperv1.ListItemRequestV1{
		Since: 0,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (k *KeeperClient) SetItem(ctx context.Context, email string, password string) (string, error) {
	const op = "client.keeper.Login"

	resp, err := k.api.SetItemV1((ctx, &keeperv1.SetItemRequestV1{
		Item: keeperv1.Item{

		}
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Token, nil
}
