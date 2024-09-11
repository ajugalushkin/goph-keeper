package minio

import (
	"github.com/ajugalushkin/goph-keeper/server/config"
	"testing"
)

// Returns an error if the Minio client cannot be created due to invalid credentials
func TestNewMinioStorage_InvalidCredentials(t *testing.T) {
	cfg := config.Minio{
		Endpoint: "localhost:9000",
		Username: "invaliduser",
		Password: "invalidpass",
		SSL:      false,
		Bucket:   "testbucket",
	}

	storage, err := NewMinioStorage(cfg)
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	if storage != nil {
		t.Fatalf("expected storage to be nil")
	}
}
