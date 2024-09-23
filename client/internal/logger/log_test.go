package logger

import (
	"testing"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"

	"github.com/stretchr/testify/assert"
)

// Returns a LogInstance when called for the first time
func TestReturnsLogInstanceFirstCall(t *testing.T) {
	config.GetConfig().Env = "dev"
	logInstance := GetLogger()

	assert.NotNil(t, logInstance)
	assert.NotNil(t, logInstance.Log)
}

// Handles missing or invalid environment configuration gracefully
func TestHandlesInvalidEnvConfigGracefully(t *testing.T) {
	config.GetConfig().Env = "invalid_env"
	logInstance := GetLogger()

	assert.NotNil(t, logInstance)
	assert.NotNil(t, logInstance.Log)
}
