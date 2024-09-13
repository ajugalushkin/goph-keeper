package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MustLoad returns a valid Config object when a valid config file is provided
func TestMustLoadValidConfig(t *testing.T) {
	// Set up a temporary config file
	configContent := `
    env: "development"
    storage:
      path: "/tmp/storage"
    grpc:
      address: "localhost:50051"
      timeout: "1h"
    token:
      ttl: "24h"
      secret: "supersecret"
    minio:
      endpoint: "localhost:9000"
      username: "minioadmin"
      password: "minioadmin"
      ssl: false
      bucket: "testbucket"
    `
	tmpFile, err := os.CreateTemp("", "config*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(configContent))
	assert.NoError(t, err)
	tmpFile.Close()

	// Set the CONFIG environment variable to the temp file path
	os.Setenv("SERVER_CONFIG", tmpFile.Name())
	defer os.Unsetenv("SERVER_CONFIG")

	// Load the config
	cfg := MustLoad()

	// Assert the config values
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, "/tmp/storage", cfg.Storage.Path)
	assert.Equal(t, "localhost:50051", cfg.GRPC.Address)
	assert.Equal(t, 1*time.Hour, cfg.GRPC.Timeout)
	assert.Equal(t, 24*time.Hour, cfg.Token.TTL)
	assert.Equal(t, "supersecret", cfg.Token.Secret)
	assert.Equal(t, "localhost:9000", cfg.Minio.Endpoint)
	assert.Equal(t, "minioadmin", cfg.Minio.Username)
	assert.Equal(t, "minioadmin", cfg.Minio.Password)
	assert.False(t, cfg.Minio.SSL)
	assert.Equal(t, "testbucket", cfg.Minio.Bucket)
}

// Handle missing configuration file gracefully
func TestHandleMissingConfigFileGracefully(t *testing.T) {
	configPath := "non_existent_config.yaml"
	os.Setenv("TOKEN_SECRET", "test_secret")
	defer os.Unsetenv("TOKEN_SECRET")

	cfg := MustLoadByPath(configPath)

	assert.NotNil(t, cfg)
	assert.Equal(t, "test_secret", cfg.Token.Secret)
}
