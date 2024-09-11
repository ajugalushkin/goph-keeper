package cmd

import (
	"github.com/ajugalushkin/goph-keeper/client/internal/token"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInit_TokenStorageWithFileBasedStorageWhenFileDoesNotExist(t *testing.T) {
	// Arrange
	tokenFilePath := "nonexistent_token.txt"
	defer os.Remove(tokenFilePath) // Clean up after the test

	// Act
	TokenStorage := token.NewFileStorage(tokenFilePath)

	// Assert
	assert.NotNil(t, TokenStorage, "Token storage should not be nil")
	assert.IsType(t, &token.FileStorage{}, TokenStorage, "Token storage should be of type FileStorage")
}
func TestInit_TokenStorageWithFileBasedStorageWhenFileIsEmpty(t *testing.T) {
	// Arrange
	tokenFilePath := "empty_token.txt"
	defer os.Remove(tokenFilePath) // Clean up after the test

	// Create an empty file for testing
	file, err := os.Create(tokenFilePath)
	assert.Nil(t, err, "Error creating test file")
	defer file.Close()

	// Act
	TokenStorage := token.NewFileStorage(tokenFilePath)

	// Assert
	assert.NotNil(t, TokenStorage, "Token storage should not be nil")
	assert.IsType(t, &token.FileStorage{}, TokenStorage, "Token storage should be of type FileStorage")
}
