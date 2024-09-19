package token_cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestLoad_WhenFileContainsValidData_ShouldReturnCorrectContent(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_file_*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	expectedContent := "valid_access_token"
	_, err = tmpFile.WriteString(expectedContent)
	require.NoError(t, err)

	storage := &FileStorage{Path: tmpFile.Name()}

	actualContent, err := storage.Load()
	require.NoError(t, err)
	assert.Equal(t, expectedContent, actualContent)
}

func TestLoad_WhenFileDoesNotExist_ShouldReturnError(t *testing.T) {
	storage := &FileStorage{Path: "nonexistent_file.txt"}

	_, err := storage.Load()
	require.Error(t, err)
}

func TestLoad_WhenFileIsEmpty_ShouldReturnError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_file_*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	storage := &FileStorage{Path: tmpFile.Name()}

	_, err = storage.Load()
	require.Error(t, err)
	assert.Equal(t, io.EOF, err)
}
