package vaulttypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Returns vaultTypeBin when called on a Bin instance
func TestBinTypeReturnsVaultTypeBin(t *testing.T) {
	bin := Bin{
		FileName: "example.bin",
		Size:     1024,
	}

	result := bin.Type()

	assert.Equal(t, VaultType("bin"), result)
}

// Handles an uninitialized Bin instance gracefully
func TestUninitializedBinType(t *testing.T) {
	var bin Bin

	result := bin.Type()

	assert.Equal(t, VaultType("bin"), result)
}

// Returns the string "BINARY DATA" for any Bin instance
func TestBinStringReturnsBinaryData(t *testing.T) {
	bin := Bin{FileName: "example.txt", Size: 1024}
	expected := "BINARY DATA"
	result := bin.String()

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

// Handles empty Bin struct without errors
func TestBinStringHandlesEmptyStruct(t *testing.T) {
	bin := Bin{}
	expected := "BINARY DATA"
	result := bin.String()

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
