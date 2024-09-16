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
