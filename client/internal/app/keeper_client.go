package app

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
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
	const op = "client.keeper.Register"

	resp, err := k.api.CreateItemV1(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (k *KeeperClient) CreateItemStream(log *slog.Logger, ctx context.Context, fileName string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error("cannot open file: ", err)
		return err
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := k.api.CreateItemStreamV1(ctx)
	if err != nil {
		log.Error("cannot upload file: ", err)
		return err
	}

	req := &keeperv1.CreateItemStreamRequestV1{
		Data: &keeperv1.CreateItemStreamRequestV1_Info{
			Info: &keeperv1.CreateItemStreamRequestV1_FileInfo{
				Name: fileName,
				Type: "",
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Error("cannot send file info to server: ", err, stream.RecvMsg(nil))
		return err
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, defaultChunkSize)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error("cannot read chunk to buffer: ", err)
		}

		req := &keeperv1.CreateItemStreamRequestV1{
			Data: &keeperv1.CreateItemStreamRequestV1_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Error("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Error("cannot receive response: ", err)
		return err
	}

	log.Info("image uploaded with id: %s, size: %d", res.GetName(), res.GetSize())
	return nil
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

func GetKeeperConnection(token string) *grpc.ClientConn {
	const op = "app.GetKeeperConnection"
	log := logger.GetInstance().Log.With("op", op)

	interceptor, err := NewAuthInterceptor(token, authMethods())
	if err != nil {
		log.Error("Unable to create interceptor", "error", err)
	}

	cfg := config.GetInstance().Config
	keeperClientConnection, err := grpc.DialContext(
		context.Background(),
		cfg.Client.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Error("Unable to connect to server", "error", err)
	}

	return keeperClientConnection
}
func authMethods() map[string]bool {
	return map[string]bool{
		keeperv1.KeeperServiceV1_ListItemV1_FullMethodName: true,
		keeperv1.KeeperServiceV1_SetItemV1_FullMethodName:  true,
	}
}
