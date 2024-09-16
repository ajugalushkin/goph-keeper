package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Handle missing configuration file gracefully
func TestHandleMissingConfigFileGracefully(t *testing.T) {
	configPath := "non_existent_config.yaml"
	os.Setenv("TOKEN_SECRET", "test_secret")
	defer os.Unsetenv("TOKEN_SECRET")

	cfg := MustLoadByPath(configPath)

	assert.NotNil(t, cfg)
	assert.Equal(t, "test_secret", cfg.Token.Secret)
}

func TestFetchConfigWithCommandLineArgument(t *testing.T) {
	os.Args = []string{"cmd", "-config", "test_config.yaml"}
	expected := "test_config.yaml"
	actual := fetchConfig()
	assert.Equal(t, expected, actual)
}

// MustLoad returns a valid Config struct when a valid config file path is provided
func TestMustLoadValidConfig(t *testing.T) {
	os.Setenv("SERVER_CONFIG", "./testdata/valid_config.yaml")
	defer os.Unsetenv("SERVER_CONFIG")

	cfg := MustLoad()

	assert.NotNil(t, cfg)
	assert.Equal(t, "prod", cfg.Env)
	assert.Equal(t, "localhost:50051", cfg.GRPC.Address)
	assert.Equal(t, "mysecret", cfg.Token.Secret)
}
