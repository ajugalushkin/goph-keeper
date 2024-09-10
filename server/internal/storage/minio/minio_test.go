package minio

import (
	"testing"

	"github.com/ajugalushkin/goph-keeper/server/config"
)

// Handles invalid Minio endpoint gracefully
func TestNewMinioStorage_InvalidEndpoint(t *testing.T) {
	cfg := config.Minio{
		Endpoint: "invalid-endpoint",
		Username: "minioadmin",
		Password: "minioadmin",
		SSL:      false,
		Bucket:   "testbucket",
	}

	storage, err := NewMinioStorage(cfg)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if storage != nil {
		t.Fatalf("expected storage to be nil, got %v", storage)
	}
}
