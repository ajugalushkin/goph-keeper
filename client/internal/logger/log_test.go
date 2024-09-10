package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ajugalushkin/goph-keeper/client/config"
)

// Returns a LogInstance when called for the first time
func TestReturnsLogInstanceFirstCall(t *testing.T) {
	config.GetInstance().Config.Env = "dev"
	logInstance := GetInstance()

	assert.NotNil(t, logInstance)
	assert.NotNil(t, logInstance.Log)
}

// Handles missing or invalid environment configuration gracefully
func TestHandlesInvalidEnvConfigGracefully(t *testing.T) {
	config.GetInstance().Config.Env = "invalid_env"
	logInstance := GetInstance()

	assert.NotNil(t, logInstance)
	assert.NotNil(t, logInstance.Log)
}
