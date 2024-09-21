package vaulttypes

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// Decodes valid creds data correctly
func TestDecodeVaultWithValidCredentials(t *testing.T) {
	data := []byte(`{"type":"creds","data":{"Login":"user","Password":"pass"}}`)
	vault, err := Deserialise(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	credentials, ok := vault.(Credentials)
	if !ok {
		t.Fatalf("expected Credentials type, got %T", vault)
	}
	if credentials.Login != "user" || credentials.Password != "pass" {
		t.Errorf("expected Login: user, Password: pass, got Login: %s, Password: %s", credentials.Login, credentials.Password)
	}
}

// Handles empty input data gracefully
func TestDecodeVaultWithEmptyData(t *testing.T) {
	data := []byte(``)
	_, err := Deserialise(data)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeVaultWithUnknownSecretType(t *testing.T) {
	data := []byte(`{"type":"unknown","data":"unknown data"}`)
	_, err := Deserialise(data)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expectedError := "unknown secret type"
	if err.Error() != expectedError {
		t.Errorf("expected error: %s, got error: %s", expectedError, err.Error())
	}
}

func TestDecodeVaultWithValidTextData(t *testing.T) {
	text := Text{Data: "This is some text data"}
	content, err := Serialise(text)
	require.NoError(t, err)

	vault, err := Deserialise(content)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	text, ok := vault.(Text)
	if !ok {
		t.Fatalf("expected Text type, got %T", vault)
	}
	expectedText := "TextData: This is some text data"
	if text.String() != expectedText {
		t.Errorf("expected Text: %s, got Text: %s", expectedText, text.String())
	}
}

func TestDecodeVaultWithValidBinaryData(t *testing.T) {
	bin := Bin{
		FileName: "TestBin",
		Size:     50,
	}

	data, err := Serialise(bin)
	require.NoError(t, err)

	vault, err := Deserialise(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	bin, ok := vault.(Bin)
	if !ok {
		t.Fatalf("expected Bin type, got %T", vault)
	}
	expectedBin := "BINARY DATA"
	if bin.String() != expectedBin {
		t.Errorf("expected Bin: %s, got Bin: %s", expectedBin, bin.String())
	}
}
func TestDecodeVaultWithValidCardData(t *testing.T) {
	card := Card{
		Number:       "1234567890123456",
		ExpiryDate:   "12/25",
		SecurityCode: "123",
		Holder:       "John Doe",
	}

	data, err := Serialise(card)
	require.NoError(t, err)

	vault, err := Deserialise(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	card, ok := vault.(Card)
	if !ok {
		t.Fatalf("expected Card type, got %T", vault)
	}
	if card.Number != "1234567890123456" || card.ExpiryDate != "12/25" || card.SecurityCode != "123" {
		t.Errorf(
			"expected Number: 1234567890123456, Expiry: 12/25, CVV: 123, got Number: %s, Expiry: %s, CVV: %s",
			card.Number,
			card.ExpiryDate,
			card.SecurityCode,
		)
	}
}
