package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestContextError_Canceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := contextError(ctx)
	expectedErr := status.Error(codes.Canceled, "request is canceled")

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error: %v, got: %v", expectedErr, err)
	}
}

func TestContextError_NonCanceledNonDeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := contextError(ctx)
	require.Nil(t, err)
}
