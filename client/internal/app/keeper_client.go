package app

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"log/slog"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ajugalushkin/goph-keeper/client/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

type KeeperClient struct {
	api keeperv1.KeeperServiceV1Client
}

const defaultChunkSize = 1024 * 1024

// NewKeeperClient returns a new keeper client
func NewKeeperClient(cc *grpc.ClientConn) *KeeperClient {
	service := keeperv1.NewKeeperServiceV1Client(cc)
	return &KeeperClient{service}
}

func (k *KeeperClient) CreateItem(ctx context.Context, item *keeperv1.CreateItemRequestV1) (*keeperv1.CreateItemResponseV1, error) {
	const op = "client.keeper.CreateItem"

	resp, err := k.api.CreateItemV1(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (k *KeeperClient) CreateItemStream(log *slog.Logger, ctx context.Context, fileName string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error("cannot open file: ", slog.String("error", err.Error()))
		return err
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := k.api.CreateItemStreamV1(ctx)
	if err != nil {
		log.Error("cannot upload file: ", slog.String("error", err.Error()))
		return err
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, defaultChunkSize)

	detectReader, err := mimetype.DetectReader(reader)
	if err != nil {
		log.Error("cannot detect reader: ", slog.String("error", err.Error()))
	}

	req := &keeperv1.CreateItemStreamRequestV1{
		Data: &keeperv1.CreateItemStreamRequestV1_Info{
			Info: &keeperv1.CreateItemStreamRequestV1_FileInfo{
				Name: fileName,
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Error("cannot send file info to server: ",
			slog.String("error", err.Error()),
			slog.String("stream msg", stream.RecvMsg(nil).Error()))
		return err
	}

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error("cannot read chunk to buffer: ", slog.String("error", err.Error()))
		}

		req := &keeperv1.CreateItemStreamRequestV1{
			Data: &keeperv1.CreateItemStreamRequestV1_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Error("cannot send chunk to server:",
				slog.String("error", err.Error()),
				slog.String("stream msg", stream.RecvMsg(nil).Error()))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Error("cannot receive response: ", slog.String("error", err.Error()))
		return err
	}

	log.Info(
		"image uploaded with id: ",
		slog.String("name", res.GetName()),
		slog.String("size", strconv.Itoa(int(res.GetSize()))),
	)
	return nil
}

func (k *KeeperClient) UpdateItem(ctx context.Context, item *keeperv1.UpdateItemRequestV1) (*keeperv1.UpdateItemResponseV1, error) {
	const op = "keeper.UpdateItem"

	resp, err := k.api.UpdateItemV1(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (k *KeeperClient) DeleteItem(ctx context.Context, item *keeperv1.DeleteItemRequestV1) (*keeperv1.DeleteItemResponseV1, error) {
	const op = "keeper.DeleteItem"

	resp, err := k.api.DeleteItemV1(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (k *KeeperClient) GetItem(ctx context.Context, item *keeperv1.GetItemRequestV1) (*keeperv1.GetItemResponseV1, error) {
	const op = "client.keeper.GetItem"

	resp, err := k.api.GetItemV1(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (k *KeeperClient) ListItems(ctx context.Context, item *keeperv1.ListItemsRequestV1) (*keeperv1.ListItemsResponseV1, error) {
	const op = "client.keeper.Register"

	list, err := k.api.ListItemsV1(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func GetKeeperConnection(token string) *grpc.ClientConn {
	const op = "app.GetKeeperConnection"
	log := logger.GetInstance().Log.With("op", op)

	interceptor, err := NewAuthInterceptor(token, authMethods())
	if err != nil {
		log.Error("Unable to create interceptor: ", slog.String("error", err.Error()))
	}

	cfg := config.GetInstance().Config
	keeperClientConnection, err := grpc.NewClient(
		cfg.Client.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Error("Unable to connect to server: ", slog.String("error", err.Error()))
	}

	return keeperClientConnection
}
func authMethods() map[string]bool {
	return map[string]bool{
		keeperv1.KeeperServiceV1_ListItemsV1_FullMethodName:        true,
		keeperv1.KeeperServiceV1_GetItemV1_FullMethodName:          true,
		keeperv1.KeeperServiceV1_CreateItemV1_FullMethodName:       true,
		keeperv1.KeeperServiceV1_CreateItemStreamV1_FullMethodName: true,
	}
}
