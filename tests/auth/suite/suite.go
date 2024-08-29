package suite

//type Suite struct {
//	*testing.T
//	Cfg        *config.Config
//	AuthClient v1.KeeperServiceV1Client
//}
//
//func New(t *testing.T) (context.Context, *Suite) {
//	t.Helper()
//	t.Parallel()
//
//	cfg := config.MustLoadByPath("../../config/dev.yaml")
//	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
//
//	t.Cleanup(func() {
//		t.Helper()
//		cancelCtx()
//	})
//
//	cc, err := grpc.NewClient(cfg.GRPC.ServerAddress,
//		grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		t.Fatalf("grpc server connection failed: %v", err)
//	}
//
//	return ctx, &Suite{
//		T:          t,
//		Cfg:        cfg,
//		AuthClient: v1.NewKeeperServiceV1Client(cc),
//	}
//}
