package vaulttypes

//import (
//	"encoding/json"
//	"testing"
//
//	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes/mocks"
//)
//
//// EncodeVault successfully encodes a Vault object into JSON
//func TestEncodeVaultSuccess(t *testing.T) {
//	mockVault := mocks.NewVault(t)
//
//	encoded, err := EncodeVault(mockVault)
//	if err != nil {
//		t.Fatalf("expected no error, got %v", err)
//	}
//
//	var c container
//	if err := json.Unmarshal(encoded, &c); err != nil {
//		t.Fatalf("expected valid JSON, got error %v", err)
//	}
//
//	if c.Type != vaultTypeText {
//		t.Errorf("expected type %v, got %v", vaultTypeText, c.Type)
//	}
//
//	var decodedData string
//	if err := json.Unmarshal(c.Data, &decodedData); err != nil {
//		t.Fatalf("expected valid data JSON, got error %v", err)
//	}
//
//	if decodedData != "example text" {
//		t.Errorf("expected data %v, got %v", "example text", decodedData)
//	}
//}
