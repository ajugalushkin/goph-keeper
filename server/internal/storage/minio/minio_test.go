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
func TestNewMinioStorage_EmptyEndpoint(t *testing.T) {
	cfg := config.Minio{
		Endpoint: "",
		Username: "testuser",
		Password: "testpass",
		SSL:      false,
		Bucket:   "testbucket",
	}

	storage, err := NewMinioStorage(cfg)
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	expectedError := "Endpoint:  does not follow ip address or domain name standards."
	if err.Error() != expectedError {
		t.Fatalf("expected error: %s, got: %s", expectedError, err.Error())
	}

	if storage != nil {
		t.Fatalf("expected storage to be nil")
	}
}
