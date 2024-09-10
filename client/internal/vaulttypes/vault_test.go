package vaulttypes

import (
	"testing"
)

type mockVault struct {
	vaultType VaultType
	data      string
}

func (m mockVault) Type() VaultType {
	return m.vaultType
}

func (m mockVault) String() string {
	return m.data
}

// Decodes valid credentials data correctly
func TestDecodeVaultWithValidCredentials(t *testing.T) {
	data := []byte(`{"type":"credentials","data":{"Login":"user","Password":"pass"}}`)
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
