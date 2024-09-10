package token

import (
	"os"
	"testing"
)

// Creates a new FileStorage instance with the specified path
func TestCreatesNewFileStorageWithPath(t *testing.T) {
	path := "/some/path/to/token"
	fs := NewFileStorage(path)

	if fs.Path != path {
		t.Errorf("expected path %s, got %s", path, fs.Path)
	}
}

// Handles empty string as path
func TestHandlesEmptyStringAsPath(t *testing.T) {
	path := ""
	fs := NewFileStorage(path)

	if fs.Path != path {
		t.Errorf("expected empty path, got %s", fs.Path)
	}
}

// Save token to a new file successfully
func TestSaveTokenSuccessfully(t *testing.T) {
	path := "test_token.txt"
	storage := NewFileStorage(path)
	defer os.Remove(path)

	err := storage.Save("test_access_token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected no error reading file, got %v", err)
	}

	if string(data) != "test_access_token" {
		t.Fatalf("expected 'test_access_token', got %s", string(data))
	}
}

// Handle error when file path is invalid
func TestSaveTokenInvalidPath(t *testing.T) {
	invalidPath := "/invalid_path/test_token.txt"
	storage := NewFileStorage(invalidPath)

	err := storage.Save("test_access_token")
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

// Successfully loads token from a valid file path
func TestLoadTokenFromValidFilePath(t *testing.T) {
	// Arrange
	path := "test_token.txt"
	expectedToken := "test_token"
	err := os.WriteFile(path, []byte(expectedToken), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	defer os.Remove(path)

	storage := NewFileStorage(path)

	// Act
	token, err := storage.Load()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if token != expectedToken {
		t.Errorf("Expected token %v, got %v", expectedToken, token)
	}
}

// File does not exist at the specified path
func TestLoadTokenFromNonExistentFilePath(t *testing.T) {
	// Arrange
	path := "non_existent_token.txt"
	storage := NewFileStorage(path)

	// Act
	token, err := storage.Load()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if token != "" {
		t.Errorf("Expected empty token, got %v", token)
	}
}
