package suite

import (
	"context"
	"log/slog"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc/test/bufconn"

	"github.com/ajugalushkin/goph-keeper/server/config"
	authhandlerv1 "github.com/ajugalushkin/goph-keeper/server/internal/handlers/grpc/auth/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/services"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage/postgres"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient v1.AuthServiceV1Client
	Closer     func()
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

	baseServer := grpc.NewServer()

	userStorage, err := postgres.NewUserStorage(cfg.Storage.Path)
	if err != nil {
		panic(err)
	}

	jwtManager := services.NewJWTManager(log, cfg.Token.Secret, cfg.Token.TTL)

	serviceAuth := services.NewAuthService(log, userStorage, userStorage, jwtManager)
	authhandlerv1.Register(baseServer, serviceAuth)

	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Error("error serving server: %v", err)
		}
	}()

	cc, err := grpc.NewClient(cfg.GRPC.Address,
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Error("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: v1.NewAuthServiceV1Client(cc),
		Closer:     closer,
	}
}
