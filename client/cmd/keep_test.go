package cmd

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
)

// keepCmd is successfully added to rootCmd
func TestKeepCmdAddedToRootCmd(t *testing.T) {
	rootCmd := &cobra.Command{}
	keepCmd := &cobra.Command{Use: "keep"}

	rootCmd.AddCommand(keepCmd)

	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "keep" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("keepCmd was not added to rootCmd")
	}
}

// keepCmd is nil when added to rootCmd
func TestNilKeepCmdAddedToRootCmd(t *testing.T) {
	rootCmd := &cobra.Command{}
	var keepCmd *cobra.Command = nil

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when adding nil keepCmd, but did not panic")
		}
	}()

	rootCmd.AddCommand(keepCmd)
}

// Encrypts a valid Vault object successfully
func TestEncryptSecretSuccess(t *testing.T) {
	vault := vaulttypes.Text{Data: "example text"}
	encrypted, err := encryptSecret(vault)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(encrypted) == 0 {
		t.Fatalf("expected non-empty encrypted data")
	}
}

// Successfully decrypt valid byte slice into Vault object
func TestDecryptSecretValidByteSlice(t *testing.T) {
	// Arrange
	vault := vaulttypes.Text{Data: "example"}
	encodedVault, err := vaulttypes.EncodeVault(vault)
	if err != nil {
		t.Fatalf("Failed to encode vault: %v", err)
	}

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	data := Data{Context: encodedVault}
	if err := enc.Encode(data); err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}

	// Act
	result, err := decryptSecret(buff.Bytes())

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.String() != vault.String() {
		t.Errorf("Expected %v, got %v", vault.String(), result.String())
	}
}

// Handle empty byte slice input
func TestDecryptSecretEmptyByteSlice(t *testing.T) {
	// Arrange
	emptyBytes := []byte{}

	// Act
	_, err := decryptSecret(emptyBytes)

	// Assert
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}
