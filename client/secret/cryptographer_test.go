package secret

import (
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCryptographer_Encrypt(t *testing.T) {
	c := NewCryptographer()
	card := vaulttypes.Card{
		Number:       "123",
		ExpiryDate:   "01/2022",
		SecurityCode: "123",
		Holder:       "John Doe",
	}

	encryptedData, err := c.Encrypt(card)
	require.NoError(t, err)
	require.NotNil(t, encryptedData)

	decryptedVault, err := c.Decrypt(encryptedData)
	require.NoError(t, err)
	require.Equal(t, card, decryptedVault)
}

func TestCryptographer_Encrypt_EmptyVault(t *testing.T) {
	c := NewCryptographer()

	emptyVault := vaulttypes.Card{
		Number:       "123",
		ExpiryDate:   "01/2022",
		SecurityCode: "123",
		Holder:       "John Doe",
	}
	encryptedData, err := c.Encrypt(emptyVault)

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if len(encryptedData) == 0 {
		t.Error("Expected non-empty encrypted data, but got empty data")
	}
}
