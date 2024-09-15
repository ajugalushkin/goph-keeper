package suite

import (
	"context"
	"testing"

	"github.com/ajugalushkin/goph-keeper/server/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/ajugalushkin/goph-keeper/gen/auth/v1"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient v1.AuthServiceV1Client
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	//t.Parallel()

	cfg := config.MustLoadByPath("../../server/config/config.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(cfg.GRPC.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: v1.NewAuthServiceV1Client(cc),
	}
}
