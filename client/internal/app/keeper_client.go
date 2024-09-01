package app

import (
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"google.golang.org/grpc"
)

type KeeperClient struct {
	service keeperv1.KeeperServiceV1Client
}

// NewKeeperClient returns a new keeper client
func NewKeeperClient(cc *grpc.ClientConn) *KeeperClient {
	service := keeperv1.NewKeeperServiceV1Client(cc)
	return &KeeperClient{service}
}
