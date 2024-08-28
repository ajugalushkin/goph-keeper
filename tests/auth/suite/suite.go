package suite

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ajugalushkin/goph-keeper/internal/config"

	authv1 "github.com/ajugalushkin/goph-keeper/common/gen/keeper/v1"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient authv1.AuthV1Client
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../../config/server_dev.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(cfg.GRPC.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: authv1.NewAuthV1Client(cc),
	}
}
