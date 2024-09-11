package minio

import (
	"github.com/ajugalushkin/goph-keeper/server/config"
	"strings"
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
func TestNewMinioStorage_EmptyUsername(t *testing.T) {
	cfg := config.Minio{
		Endpoint: "localhost:9000",
		Username: "",
		Password: "validpass",
		SSL:      false,
		Bucket:   "testbucket",
	}

	_, err := NewMinioStorage(cfg)
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	expectedError := "Get \"http://localhost:9000/testbucket/?location=\": dial tcp [::1]:9000: connectex: No connection could be made because the target machine actively refused it."
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}
func TestNewMinioStorage_EmptyPassword(t *testing.T) {
	cfg := config.Minio{
		Endpoint: "localhost:9000",
		Username: "testuser",
		Password: "",
		SSL:      false,
		Bucket:   "testbucket",
	}

	_, err := NewMinioStorage(cfg)
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	expectedError := "Get \"http://localhost:9000/testbucket/?location=\": dial tcp [::1]:9000: connectex: No connection could be made because the target machine actively refused it."
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}
