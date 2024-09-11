package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

//go:generate mockery --name CreateItemStreamClient
type CreateItemStreamClient interface {
	keeperv1.KeeperServiceV1_CreateItemStreamV1Client
}

type KeeperClient struct {
	api keeperv1.KeeperServiceV1Client
}

// NewKeeperClient returns a new keeper client
func NewKeeperClient(cc *grpc.ClientConn) *KeeperClient {
	service := keeperv1.NewKeeperServiceV1Client(cc)
	return &KeeperClient{service}
}

func (k *KeeperClient) CreateItem(
	ctx context.Context,
	item *keeperv1.CreateItemRequestV1,
) (*keeperv1.CreateItemResponseV1, error) {
	return k.api.CreateItemV1(ctx, item)
}

func (k *KeeperClient) CreateItemStream(
	ctx context.Context,
	name string,
	file *os.File,
	content []byte,
) (*keeperv1.CreateItemResponseV1, error) {
	stream, err := k.api.CreateItemStreamV1(context.Background())
	if err != nil {
		slog.Error("cannot upload file: ", slog.String("error", err.Error()))
		return nil, err
	}

	req := &keeperv1.CreateItemStreamRequestV1{
		Data: &keeperv1.CreateItemStreamRequestV1_Info{
			Info: &keeperv1.CreateItemStreamRequestV1_FileInfo{
				Name:    name,
				Content: content,
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		slog.Error("cannot send file info to server: ",
			slog.String("error", err.Error()),
			slog.String("stream msg", stream.RecvMsg(nil).Error()))
		return nil, err
	}

	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			slog.Error("cannot read chunk to buffer: ", slog.String("error", err.Error()))
		}

		req := &keeperv1.CreateItemStreamRequestV1{
			Data: &keeperv1.CreateItemStreamRequestV1_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			slog.Error("cannot send chunk to server:",
				slog.String("error", err.Error()),
				slog.String("stream msg", stream.RecvMsg(nil).Error()))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		slog.Error("cannot receive response: ", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info(
		"file uploaded with id: ",
		slog.String("name", res.GetName()),
		slog.String("size", res.GetVersion()),
	)
	return res, nil
}

func (k *KeeperClient) UpdateItem(
	ctx context.Context,
	item *keeperv1.UpdateItemRequestV1,
) (*keeperv1.UpdateItemResponseV1, error) {
	const op = "keeper.UpdateItem"

	resp, err := k.api.UpdateItemV1(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (k *KeeperClient) DeleteItem(
	ctx context.Context,
	item *keeperv1.DeleteItemRequestV1,
) (*keeperv1.DeleteItemResponseV1, error) {
	const op = "keeper.DeleteItem"

	resp, err := k.api.DeleteItemV1(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (k *KeeperClient) GetItem(
	ctx context.Context,
	item *keeperv1.GetItemRequestV1,
) (*keeperv1.GetItemResponseV1, error) {
	const op = "client.keeper.GetItem"

	resp, err := k.api.GetItemV1(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (k *KeeperClient) GetFile(
	ctx context.Context,
	name string,
) (keeperv1.KeeperServiceV1_GetItemStreamV1Client, error) {
	const op = "client.keeper.GetItem"

	stream, err := k.api.GetItemStreamV1(
		ctx,
		&keeperv1.GetItemRequestV1{Name: name},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stream, nil
}

func (k *KeeperClient) ListItems(
	ctx context.Context,
	item *keeperv1.ListItemsRequestV1,
) (*keeperv1.ListItemsResponseV1, error) {
	const op = "client.keeper.Register"

	list, err := k.api.ListItemsV1(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func GetKeeperConnection(
	log *slog.Logger,
	address string,
	token string,
) *grpc.ClientConn {
	const op = "app.GetKeeperConnection"
	log.With("op", op)

	interceptor, err := NewAuthInterceptor(token, authMethods())
	if err != nil {
		log.Error("Unable to create interceptor: ", slog.String("error", err.Error()))
	}

	if interceptor == nil {
		log.Error("interceptor is nil")
		return nil
	}

	keeperClientConnection, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Error(
			"Unable to connect to server: ",
			slog.String("error", err.Error()),
		)
	}

	return keeperClientConnection
}
func authMethods() map[string]bool {
	return map[string]bool{
		keeperv1.KeeperServiceV1_ListItemsV1_FullMethodName:        true,
		keeperv1.KeeperServiceV1_GetItemV1_FullMethodName:          true,
		keeperv1.KeeperServiceV1_CreateItemV1_FullMethodName:       true,
		keeperv1.KeeperServiceV1_CreateItemStreamV1_FullMethodName: true,
		keeperv1.KeeperServiceV1_DeleteItemV1_FullMethodName:       true,
		keeperv1.KeeperServiceV1_UpdateItemV1_FullMethodName:       true,
	}
}
