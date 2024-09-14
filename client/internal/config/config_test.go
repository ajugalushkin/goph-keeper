package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// GetInstance returns a singleton instance of CfgInstance
func TestGetInstanceReturnsSingleton(t *testing.T) {
	instance1 := GetConfig()
	instance2 := GetConfig()

	assert.NotNil(t, instance1)
	assert.Equal(t, instance1, instance2)
}
