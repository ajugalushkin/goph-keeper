package config

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Mkdir("test", 0777)
	newConfig := Config{
		Env:     "dev",
		Storage: Storage{Path: "postgresql://praktikum:pass@localhost:5432/goph_keeper?sslmode=disable"},
		GRPC: GRPC{
			Address: "localhost:50051",
			Timeout: time.Hour,
		},
		Token: Token{
			TTL:    time.Hour,
			Secret: "secret",
		},
		Minio: Minio{
			Endpoint: "localhost:9000",
			Username: "testuser",
			Password: "testpassword",
			SSL:      false,
			Bucket:   "testbucket",
		},
	}
	yamlData, err := yaml.Marshal(&newConfig)
	if err != nil {
		return
	}
	filePath := "test/config.yaml"
	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return
	}
	os.Setenv("SERVER_CONFIG", filePath)

	exitcode := m.Run()

	os.RemoveAll("test")
	os.Clearenv()
	os.Exit(exitcode)
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

func TestFetchConfigWithCommandLineArgument(t *testing.T) {
	os.Args = []string{"cmd", "-config", "test_config.yaml"}
	expected := "test_config.yaml"
	actual := fetchConfig()
	assert.Equal(t, expected, actual)
}

func TestMustLoadByPath_NonDefaultPath(t *testing.T) {
	// Arrange
	nonDefaultPath := "test/non_default_config.yaml"
	expectedConfig := &Config{
		Env:     "test",
		Storage: Storage{Path: "test_storage_path"},
		// Add other expected fields as per the Config struct
	}
	yamlData, err := yaml.Marshal(expectedConfig)
	require.NoError(t, err)
	err = os.WriteFile(nonDefaultPath, yamlData, 0644)
	require.NoError(t, err)

	// Act
	actualConfig := MustLoadByPath(nonDefaultPath)

	// Assert
	assert.Equal(t, expectedConfig, actualConfig)

	// Clean up
	os.Remove(nonDefaultPath)
}

func TestMustLoadByPath_MissingConfigFile(t *testing.T) {
	nonExistentConfigPath := "non_existent_config.yaml"
	expectedErrorMessage := "\"Config file not found in \" file=non_existent_config.yaml\n"

	logOutput := captureLogOutput(func() {
		MustLoadByPath(nonExistentConfigPath)
	})

	assert.Contains(t, logOutput, expectedErrorMessage)
}

func captureLogOutput(f func()) string {
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))
	slog.SetDefault(logger)

	f()

	return buf.String()
}
