package vaulttypes

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Returns vaultTypeText for any instance of Text
func TestTypeReturnsVaultTypeText(t *testing.T) {
	text := Text{Data: "example"}
	if text.Type() != VaultType("text") {
		t.Errorf("expected %v, got %v", VaultType("text"), text.Type())
	}
}

// Handles empty Text instances without errors
func TestTypeHandlesEmptyText(t *testing.T) {
	text := Text{Data: ""}
	if text.Type() != VaultType("text") {
		t.Errorf("expected %v, got %v", VaultType("text"), text.Type())
	}
}

// Returns formatted string with Data field content
func TestTextStringReturnsFormattedString(t *testing.T) {
	text := Text{Data: "example"}
	expected := "TextData: example"
	result := text.String()
	assert.Equal(t, expected, result)
}

// Handles very long Data strings without errors
func TestStringHandlesLongData(t *testing.T) {
	longData := strings.Repeat("a", 10000)
	text := Text{Data: longData}
	expected := "TextData: " + longData
	result := text.String()
	assert.Equal(t, expected, result)
}
