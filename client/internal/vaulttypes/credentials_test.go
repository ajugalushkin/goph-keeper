package vaulttypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Returns vaultTypeCredentials for any Credentials instance
func TestTypeReturnsVaultTypeCredentials(t *testing.T) {
	creds := Credentials{
		Login:    "user",
		Password: "pass",
	}

	result := creds.Type()

	assert.Equal(t, VaultType("credentials"), result)
}

// Handles empty Credentials struct without errors
func TestTypeHandlesEmptyCredentials(t *testing.T) {
	creds := Credentials{}

	result := creds.Type()

	assert.Equal(t, VaultType("credentials"), result)
}

// Returns formatted string with login and password
func TestStringReturnsFormattedString(t *testing.T) {
	creds := Credentials{
		Login:    "user123",
		Password: "pass123",
	}

	expected := "Login: user123, Password: pass123"
	result := creds.String()

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

// Handles empty login and password
func TestStringHandlesEmptyLoginAndPassword(t *testing.T) {
	creds := Credentials{
		Login:    "",
		Password: "",
	}

	expected := "Login: , Password: "
	result := creds.String()

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
