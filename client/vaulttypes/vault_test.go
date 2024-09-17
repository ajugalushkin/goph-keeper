package vaulttypes

import (
	"testing"
)

// Decodes valid creds data correctly
func TestDecodeVaultWithValidCredentials(t *testing.T) {
	data := []byte(`{"type":"creds","data":{"Login":"user","Password":"pass"}}`)
	vault, err := DecodeVault(data)
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
	_, err := DecodeVault(data)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeVaultWithUnknownSecretType(t *testing.T) {
	data := []byte(`{"type":"unknown","data":"unknown data"}`)
	_, err := DecodeVault(data)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expectedError := "unknown secret type"
	if err.Error() != expectedError {
		t.Errorf("expected error: %s, got error: %s", expectedError, err.Error())
	}
}
