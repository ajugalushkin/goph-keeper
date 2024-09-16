package vaulttypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Returns vaultTypeCard when called on a Card instance
func TestCardTypeReturnsVaultTypeCard(t *testing.T) {
	card := Card{
		Number:       "1234567890123456",
		ExpiryDate:   "12/24",
		SecurityCode: "123",
		Holder:       "John Doe",
	}

	result := card.Type()

	assert.Equal(t, VaultType("card"), result)
}

// Handles empty Card struct without errors
func TestCardTypeHandlesEmptyCardStruct(t *testing.T) {
	card := Card{}

	result := card.Type()

	assert.Equal(t, VaultType("card"), result)
}

// Correctly formats and returns a string with all card details
func TestStringMethodFormatsCardDetails(t *testing.T) {
	card := Card{
		Number:       "1234567890123456",
		ExpiryDate:   "12/24",
		SecurityCode: "123",
		Holder:       "John Doe",
	}

	expected := "Number: 1234567890123456, ExpiryDate: 12/24, SecurityCode: 123, Holder: John Doe"
	result := card.String()

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

// Handles empty strings for all card fields
func TestStringMethodHandlesEmptyFields(t *testing.T) {
	card := Card{
		Number:       "",
		ExpiryDate:   "",
		SecurityCode: "",
		Holder:       "",
	}

	expected := "Number: , ExpiryDate: , SecurityCode: , Holder: "
	result := card.String()

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
