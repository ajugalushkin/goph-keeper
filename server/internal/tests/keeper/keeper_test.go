package keeper

import (
	"context"
	"io"
	"log/slog"
	"net"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/ajugalushkin/goph-keeper/client/secret"
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	keephandlerv1 "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/services"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage/mocks"
	"github.com/ajugalushkin/goph-keeper/server/internal/tests/keeper/suite"
)

type Item struct {
	Name     string
	Email    string
	Password string
}

// Type возвращает тип хранимой информации
func (i Item) Type() vaulttypes.VaultType {
	var data vaulttypes.VaultType
	data = "item"
	return data
}

// String функция отображения приватной информации
func (i Item) String() string {
	return "ITEM DATA"
}

func TestCRUDItem_CRUDItem_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	nameExpected := gofakeit.Name()

	data := Item{
		Name:     nameExpected,
		Email:    gofakeit.Email(),
		Password: suite.RandomFakePassword(),
	}

	content, err := secret.EncryptSecret(data)
	require.NoError(t, err)

	resp, err := st.KeeperClient.CreateItemV1(ctx, &keeperv1.CreateItemRequestV1{
		Name:    nameExpected,
		Content: content,
	})
	require.NoError(t, err)
	assert.Equal(t, nameExpected, resp.GetName())

	respGet, err := st.KeeperClient.GetItemV1(ctx, &keeperv1.GetItemRequestV1{
		Name: nameExpected,
	})
	require.NoError(t, err)
	assert.Equal(t, nameExpected, respGet.GetName())
	assert.Equal(t, resp.GetVersion(), respGet.GetVersion())

	dataUpd := Item{
		Name:     nameExpected,
		Email:    gofakeit.Email(),
		Password: suite.RandomFakePassword(),
	}

	contentUpd, err := secret.EncryptSecret(dataUpd)
	require.NoError(t, err)

	respUpd, err := st.KeeperClient.UpdateItemV1(ctx, &keeperv1.UpdateItemRequestV1{
		Name:    nameExpected,
		Content: contentUpd,
	})
	require.NoError(t, err)
	assert.Equal(t, nameExpected, respUpd.GetName())

	respGet, err = st.KeeperClient.GetItemV1(ctx, &keeperv1.GetItemRequestV1{
		Name: nameExpected,
	})
	require.NoError(t, err)
	assert.Equal(t, nameExpected, respGet.GetName())
	assert.Equal(t, respUpd.GetVersion(), respGet.GetVersion())

	respDel, err := st.KeeperClient.DeleteItemV1(ctx, &keeperv1.DeleteItemRequestV1{
		Name: nameExpected,
	})
	require.NoError(t, err)
	assert.Equal(t, nameExpected, respDel.GetName())
}

func TestListItem_ListItem_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	for i := 0; i < 5; i++ {
		nameExpected := gofakeit.Name()

		data := Item{
			Name:     nameExpected,
			Email:    gofakeit.Email(),
			Password: suite.RandomFakePassword(),
		}

		content, err := secret.EncryptSecret(data)
		require.NoError(t, err)

		resp, err := st.KeeperClient.CreateItemV1(ctx, &keeperv1.CreateItemRequestV1{
			Name:    nameExpected,
			Content: content,
		})
		require.NoError(t, err)
		assert.Equal(t, nameExpected, resp.GetName())
	}

	_, err := st.KeeperClient.ListItemsV1(ctx, &keeperv1.ListItemsRequestV1{})
	require.NoError(t, err)
}

func TestCRUDItemStream_CRUDItem_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	fileName := "test_bin.txt"
	temp, err := os.CreateTemp("", fileName)
	require.NoError(t, err)
	defer temp.Close()

	_, err = temp.WriteString(gofakeit.Letter())
	require.NoError(t, err)

	stat, err := temp.Stat()
	require.NoError(t, err)

	fileInfo := vaulttypes.Bin{
		FileName: fileName,
		Size:     stat.Size(),
	}

	content, err := secret.EncryptSecret(fileInfo)
	require.NoError(t, err)

	req := &keeperv1.CreateItemStreamRequestV1{
		Data: &keeperv1.CreateItemStreamRequestV1_Info{
			Info: &keeperv1.CreateItemStreamRequestV1_FileInfo{
				Name:    fileName,
				Content: content,
			},
		},
	}

	stream, err := st.KeeperClient.CreateItemStreamV1(ctx)
	err = stream.Send(req)
	require.NoError(t, err)

	buffer := make([]byte, 1024)

	for {
		n, err := temp.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
		}

		req := &keeperv1.CreateItemStreamRequestV1{
			Data: &keeperv1.CreateItemStreamRequestV1_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	resp, err := stream.CloseAndRecv()
	require.NoError(t, err)

	require.NoError(t, err)
	assert.Equal(t, fileName, resp.GetName())

	streamGet, err := st.KeeperClient.GetItemStreamV1(ctx,
		&keeperv1.GetItemRequestV1{Name: fileName})
	require.NoError(t, err)

	recGet, err := streamGet.Recv()
	require.NoError(t, err)

	respSecret, err := secret.DecryptSecret(recGet.GetContent())
	require.NoError(t, err)

	fileInfoGet := respSecret.(vaulttypes.Bin)
	assert.Equal(t, fileInfo.FileName, fileInfoGet.FileName)
	assert.Equal(t, fileInfo.Size, fileInfoGet.Size)

	newFile, err := os.CreateTemp("", fileInfoGet.FileName)
	require.NoError(t, err)
	defer newFile.Close()

	for {
		req, err := streamGet.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
		}
		chunk := req.GetChunkData()

		_, err = newFile.Write(chunk)
		require.NoError(t, err)
	}

	respDel, err := st.KeeperClient.DeleteItemV1(ctx, &keeperv1.DeleteItemRequestV1{
		Name: fileName,
	})
	require.NoError(t, err)
	assert.Equal(t, fileInfo.FileName, respDel.Name)

	_, err = st.KeeperClient.GetItemStreamV1(ctx,
		&keeperv1.GetItemRequestV1{Name: fileName})
	require.NoError(t, err)
}

func TestCreateItem_CreateItem_ErrItemConflict(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	nameExpected := gofakeit.Name()

	data := Item{
		Name:     nameExpected,
		Email:    gofakeit.Email(),
		Password: suite.RandomFakePassword(),
	}

	content, err := secret.EncryptSecret(data)
	require.NoError(t, err)

	_, err = st.KeeperClient.CreateItemV1(ctx, &keeperv1.CreateItemRequestV1{
		Name:    nameExpected,
		Content: content,
	})
	require.NoError(t, err)

	_, err = st.KeeperClient.CreateItemV1(ctx, &keeperv1.CreateItemRequestV1{
		Name:    nameExpected,
		Content: content,
	})
	assert.ErrorContains(t, err, "item already exists")
}

func TestUpdateItem_UpdateItem_ErrUserNotFound(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	nameExpected := gofakeit.Name()

	data := Item{
		Name:     nameExpected,
		Email:    gofakeit.Email(),
		Password: suite.RandomFakePassword(),
	}

	content, err := secret.EncryptSecret(data)
	require.NoError(t, err)

	_, err = st.KeeperClient.UpdateItemV1(ctx, &keeperv1.UpdateItemRequestV1{
		Name:    nameExpected,
		Content: content,
	})
	require.ErrorContains(t, err, "secret not found")
}

func TestUpdateItem_UpdateItem_EmptyNameOrEmptyContent(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	nameExpected := gofakeit.Name()

	data := Item{
		Name:     nameExpected,
		Email:    gofakeit.Email(),
		Password: suite.RandomFakePassword(),
	}

	content, err := secret.EncryptSecret(data)
	require.NoError(t, err)

	_, err = st.KeeperClient.UpdateItemV1(ctx, &keeperv1.UpdateItemRequestV1{
		Name:    "",
		Content: content,
	})
	require.ErrorContains(t, err, "empty secret name")

	_, err = st.KeeperClient.UpdateItemV1(ctx, &keeperv1.UpdateItemRequestV1{
		Name:    nameExpected,
		Content: []byte(""),
	})
	require.ErrorContains(t, err, "empty secret content")
}

func TestCreateItem_CreateItem_EmptyNameOrEmptyContent(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	nameExpected := gofakeit.Name()

	data := Item{
		Name:     nameExpected,
		Email:    gofakeit.Email(),
		Password: suite.RandomFakePassword(),
	}

	content, err := secret.EncryptSecret(data)
	require.NoError(t, err)

	_, err = st.KeeperClient.CreateItemV1(ctx, &keeperv1.CreateItemRequestV1{
		Name:    "",
		Content: content,
	})
	require.ErrorContains(t, err, "empty secret name")

	_, err = st.KeeperClient.CreateItemV1(ctx, &keeperv1.CreateItemRequestV1{
		Name:    nameExpected,
		Content: []byte(""),
	})
	require.ErrorContains(t, err, "empty secret content")
}

func TestDeleteItem_DeleteItem_ErrUserNotFound(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	nameExpected := gofakeit.Name()

	_, err := st.KeeperClient.DeleteItemV1(ctx, &keeperv1.DeleteItemRequestV1{
		Name: "",
	})
	require.ErrorContains(t, err, "empty secret name")

	_, err = st.KeeperClient.DeleteItemV1(ctx, &keeperv1.DeleteItemRequestV1{
		Name: nameExpected,
	})
	require.ErrorContains(t, err, "secret not found")
}

func TestCRUDItem_EmptyUserID(t *testing.T) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	baseServer := grpc.NewServer()

	mockMinio := mocks.NewMinioStorage(t)
	vaultStorage := mocks.NewVaultStorage(t)

	serviceKeeper := services.NewKeeperService(
		log,
		vaultStorage,
		vaultStorage,
		mockMinio,
		mockMinio,
	)
	keephandlerv1.Register(baseServer, serviceKeeper)

	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Error("error serving server: ", slog.String("err", err.Error()))
		}
	}()

	cc, err := grpc.NewClient("localhost:8080",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error(
			"Unable to connect to server: ",
			slog.String("error", err.Error()),
		)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Error("error closing listener: ", slog.String("err", err.Error()))
		}
		baseServer.Stop()
	}
	defer closer()

	KeeperClient := keeperv1.NewKeeperServiceV1Client(cc)

	expectedName := gofakeit.Name()

	data := Item{
		Name:     expectedName,
		Email:    gofakeit.Email(),
		Password: suite.RandomFakePassword(),
	}

	content, err := secret.EncryptSecret(data)
	require.NoError(t, err)

	_, err = KeeperClient.CreateItemV1(context.Background(), &keeperv1.CreateItemRequestV1{
		Name:    expectedName,
		Content: content,
	})
	assert.ErrorContains(t, err, "empty user id")

	_, err = KeeperClient.UpdateItemV1(context.Background(), &keeperv1.UpdateItemRequestV1{
		Name:    expectedName,
		Content: content,
	})
	assert.ErrorContains(t, err, "empty user id")

	_, err = KeeperClient.DeleteItemV1(context.Background(), &keeperv1.DeleteItemRequestV1{
		Name: expectedName,
	})
	assert.ErrorContains(t, err, "empty user id")

	_, err = KeeperClient.GetItemV1(context.Background(), &keeperv1.GetItemRequestV1{
		Name: expectedName,
	})
	assert.ErrorContains(t, err, "empty user id")

	_, err = KeeperClient.ListItemsV1(context.Background(), &keeperv1.ListItemsRequestV1{})
	assert.ErrorContains(t, err, "empty user id")
}

func TestListItem_ListItem_ErrUserNotFound(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	_, err := st.KeeperClient.ListItemsV1(ctx, &keeperv1.ListItemsRequestV1{})
	require.ErrorContains(t, err, "failed to list secrets")
}

func TestGetItem_GetItem_ErrUserNotFound(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	nameExpected := gofakeit.Name()

	_, err := st.KeeperClient.GetItemV1(ctx, &keeperv1.GetItemRequestV1{
		Name: "",
	})
	require.ErrorContains(t, err, "secret name is empty")

	_, err = st.KeeperClient.GetItemV1(ctx, &keeperv1.GetItemRequestV1{
		Name: nameExpected,
	})
	require.ErrorContains(t, err, "failed to get item")
}
