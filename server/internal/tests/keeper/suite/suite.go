package suite

import (
	"context"
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

	clientInter "github.com/ajugalushkin/goph-keeper/client/interceptors"
	authv1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/config"
	serverInter "github.com/ajugalushkin/goph-keeper/server/interceptors"
	authhandlerv1 "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/auth/v1"
	keephandlerv1 "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/services"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage/mocks"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage/postgres"
)

type Suite struct {
	*testing.T
	Cfg          *config.Config
	KeeperClient keeperv1.KeeperServiceV1Client
	Closer       func()
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()

	cfg := config.MustLoadByPath("./suite/config.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	jwtManager := services.NewJWTManager(log, cfg.Token.Secret, cfg.Token.TTL)

	serverInterceptor := serverInter.NewAuthInterceptor(log, jwtManager, accessibleMethods())

	baseServer := grpc.NewServer(
		grpc.UnaryInterceptor(serverInterceptor.Unary()),
		grpc.StreamInterceptor(serverInterceptor.Stream()))

	userStorage, err := postgres.NewUserStorage(cfg.Storage.Path)
	if err != nil {
		panic(err)
	}

	serviceAuth := services.NewAuthService(log, userStorage, userStorage, jwtManager)
	authhandlerv1.Register(baseServer, serviceAuth)

	vaultStorage, err := postgres.NewVaultStorage(cfg.Storage.Path)
	if err != nil {
		panic(err)
	}

	minioStorage := mocks.NewMinioStorage(t)
	//minioStorage, err := minio.NewMinioStorage(cfg.Minio)
	//if err != nil {
	//	panic(err)
	//}

	serviceKeeper := services.NewKeeperService(
		log,
		vaultStorage,
		vaultStorage,
		minioStorage,
		minioStorage,
	)
	keephandlerv1.Register(baseServer, serviceKeeper)

	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Error("error serving server: ", slog.String("err", err.Error()))
		}
	}()

	authConn, err := grpc.NewClient(cfg.GRPC.Address,
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("error connecting to server: ", slog.String("err", err.Error()))
	}

	authClient := authv1.NewAuthServiceV1Client(authConn)

	email := gofakeit.Email()
	pass := RandomFakePassword()

	respReq, err := authClient.RegisterV1(ctx, &authv1.RegisterRequestV1{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReq.GetUserId())

	respLog, err := authClient.LoginV1(ctx, &authv1.LoginRequestV1{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)

	token := respLog.GetToken()
	require.NotEmpty(t, token)

	clientInterceptor, err := clientInter.NewAuthInterceptor(token, authMethods())
	if err != nil {
		log.Error("Unable to create interceptor: ", slog.String("error", err.Error()))
	}
	if clientInterceptor == nil {
		log.Error("interceptor is nil ", slog.String("err", err.Error()))
		return nil, nil
	}

	cc, err := grpc.NewClient(cfg.GRPC.Address,
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(clientInterceptor.Unary()),
		grpc.WithStreamInterceptor(clientInterceptor.Stream()),
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

	return ctx, &Suite{
		T:            t,
		Cfg:          cfg,
		KeeperClient: keeperv1.NewKeeperServiceV1Client(cc),
		Closer:       closer,
	}
}

const (
	passDefaultLen = 10
)

func RandomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
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

func accessibleMethods() []string {
	return []string{
		authv1.AuthServiceV1_RegisterV1_FullMethodName,
		authv1.AuthServiceV1_LoginV1_FullMethodName,
	}
}
