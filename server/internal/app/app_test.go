package app

import (
	"github.com/ajugalushkin/goph-keeper/server/config"
	"github.com/ajugalushkin/goph-keeper/server/internal/storage/postgres"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInitKeeperService_NilLogger(t *testing.T) {
	cfg := &config.Config{
		Env:     "dev",
		Storage: config.Storage{Path: "postgresql://praktikum:pass@localhost:5432/goph_keeper?sslmode=disable"},
		GRPC: config.GRPC{
			Address: "localhost:50051",
			Timeout: time.Hour,
		},
		Token: config.Token{
			TTL:    time.Hour,
			Secret: "secret",
		},
		Minio: config.Minio{
			Endpoint: "localhost:9000",
			Username: "testuser",
			Password: "testpassword",
			SSL:      false,
			Bucket:   "testbucket",
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("initKeeperService should not panic with nil logger: %v", r)
		}
	}()

	initKeeperService(nil, cfg)
}

func TestInitKeeperService_InvalidStoragePath(t *testing.T) {
	cfg := &config.Config{
		Env:     "dev",
		Storage: config.Storage{Path: "invalid_storage_path"},
		GRPC: config.GRPC{
			Address: "localhost:50051",
			Timeout: time.Hour,
		},
		Token: config.Token{
			TTL:    time.Hour,
			Secret: "secret",
		},
		Minio: config.Minio{
			Endpoint: "localhost:9000",
			Username: "testuser",
			Password: "testpassword",
			SSL:      false,
			Bucket:   "testbucket",
		},
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("initKeeperService should panic with invalid storage path")
		} else if err, ok := r.(error); !ok || !strings.Contains(err.Error(), "invalid_storage_path") {
			t.Errorf("initKeeperService should panic with 'invalid_storage_path' error, but got: %v", r)
		}
	}()

	initKeeperService(nil, cfg)
}

func TestInitKeeperService_InvalidMinioConfig(t *testing.T) {
	cfg := &config.Config{
		Env:     "dev",
		Storage: config.Storage{Path: "postgresql://praktikum:pass@localhost:5432/goph_keeper?sslmode=disable"},
		GRPC: config.GRPC{
			Address: "localhost:50051",
			Timeout: time.Hour,
		},
		Token: config.Token{
			TTL:    time.Hour,
			Secret: "secret",
		},
		Minio: config.Minio{
			Endpoint: "",
			Username: "testuser",
			Password: "testpassword",
			SSL:      false,
			Bucket:   "testbucket",
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("initKeeperService should panic with invalid Minio configuration")
		}
	}()

	initKeeperService(nil, cfg)
}
func TestInitKeeperService_NilVaultStorage(t *testing.T) {
	cfg := &config.Config{
		Env:     "dev",
		Storage: config.Storage{Path: "postgresql://praktikum:pass@localhost:5432/goph_keeper?sslmode=disable"},
		GRPC: config.GRPC{
			Address: "localhost:50051",
			Timeout: time.Hour,
		},
		Token: config.Token{
			TTL:    time.Hour,
			Secret: "secret",
		},
		Minio: config.Minio{
			Endpoint: "localhost:9000",
			Username: "testuser",
			Password: "testpassword",
			SSL:      false,
			Bucket:   "testbucket",
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("initKeeperService should panic with nil vault storage")
		}
	}()

	initKeeperService(slog.New(slog.NewTextHandler(os.Stdout, nil)), &config.Config{
		Storage: config.Storage{Path: ""},
		Minio:   cfg.Minio,
	})
}
func TestInitKeeperService_NilMinioStorage(t *testing.T) {
	cfg := &config.Config{
		Env:     "dev",
		Storage: config.Storage{Path: "postgresql://praktikum:pass@localhost:5432/goph_keeper?sslmode=disable"},
		GRPC: config.GRPC{
			Address: "localhost:50051",
			Timeout: time.Hour,
		},
		Token: config.Token{
			TTL:    time.Hour,
			Secret: "secret",
		},
		Minio: config.Minio{
			Endpoint: "localhost:9000",
			Username: "testuser",
			Password: "testpassword",
			SSL:      false,
			Bucket:   "testbucket",
		},
	}

	_, err := postgres.NewVaultStorage(cfg.Storage.Path)
	if err != nil {
		t.Fatalf("failed to create vault storage: %v", err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("initKeeperService should panic with nil minio storage")
		}
	}()

	initKeeperService(nil, &config.Config{
		Minio: config.Minio{}, // nil Minio configuration
	})
}
