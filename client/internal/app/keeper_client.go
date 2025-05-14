package app

import (
	"context"

	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

//go:generate mockery --name CreateItemStreamClient
type CreateItemStreamClient interface {
	keeperv1.KeeperServiceV1_CreateItemStreamV1Client
}

//go:generate mockery --name KeeperClient
type KeeperClient interface {
	CreateItem(
		ctx context.Context,
		item *keeperv1.CreateItemRequestV1,
	) (*keeperv1.CreateItemResponseV1, error)
	CreateItemStream(
		ctx context.Context,
		name string,
		filePath string,
	) (*keeperv1.CreateItemResponseV1, error)
	UpdateItem(
		ctx context.Context,
		item *keeperv1.UpdateItemRequestV1,
	) (*keeperv1.UpdateItemResponseV1, error)
	DeleteItem(
		ctx context.Context,
		item *keeperv1.DeleteItemRequestV1,
	) (*keeperv1.DeleteItemResponseV1, error)
	GetItem(
		ctx context.Context,
		item *keeperv1.GetItemRequestV1,
	) (*keeperv1.GetItemResponseV1, error)
	GetFile(
		ctx context.Context,
		name string,
		path string,
	) error
	ListItems(
		ctx context.Context,
		item *keeperv1.ListItemsRequestV1,
	) (*keeperv1.ListItemsResponseV1, error)
}
