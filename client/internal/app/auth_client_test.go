package app

import (
	"testing"

	"google.golang.org/grpc"
)

// Creates an AuthClient instance with a valid grpc.ClientConn
func TestNewAuthClientWithValidConn(t *testing.T) {
	// Arrange
	conn := new(grpc.ClientConn)

	// Act
	client := NewAuthClient(conn)

	// Assert
	if client == nil {
		t.Errorf("Expected AuthClient instance, got nil")
	}
	if client.api == nil {
		t.Errorf("Expected AuthServiceV1Client instance, got nil")
	}
}

// Handles nil grpc.ClientConn input gracefully
func TestNewAuthClientWithNilConn(t *testing.T) {
	// Arrange
	var conn *grpc.ClientConn = nil

	// Act
	client := NewAuthClient(conn)

	// Assert
	if client == nil {
		t.Errorf("Expected AuthClient instance, got nil")
	}
	if client.api == nil {
		t.Errorf("Expected AuthServiceV1Client instance, got nil")
	}
}
