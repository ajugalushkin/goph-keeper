package app

import (
	"log/slog"
	"os"
	"testing"

	"github.com/ajugalushkin/goph-keeper/server/config"
)

// Configuration paths for storage are invalid or inaccessible
func TestNewAppWithInvalidStoragePath(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	cfg := &config.Config{
		Storage: config.Storage{Path: "/invalid/path"},
		Token:   config.Token{Secret: "secret", TTL: 3600},
		GRPC:    config.GRPC{Address: "localhost:50051"},
		Minio:   config.Minio{Endpoint: "localhost:9000", Username: "user", Password: "pass", SSL: false, Bucket: "bucket"},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	_ = New(log, cfg)
}
