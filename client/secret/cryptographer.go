package secret

import (
	"bytes"
	"encoding/gob"

	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
)

//go:generate mockery --name Cipher
type Cipher interface {
	Encrypt(vaulttypes.Vault) ([]byte, error)
	Decrypt(data []byte) (vaulttypes.Vault, error)
}

// Cryptographer is a struct used for encoding and decoding Vault data using Gob encoding.
type Cryptographer struct {
	Context []byte
}

func NewCryptographer() *Cryptographer {
	return &Cryptographer{}
}

// Encrypt takes a Vault as input and encrypts it using the Gob encoding.
// The function first encodes the Vault into a byte slice using the Serialise function from the vaulttypes package.
// Then, it creates a new Cryptographer struct containing the encoded byte slice.
// The Cryptographer struct is encoded using the Gob encoding into a bytes.Buffer.
// Finally, the bytes.Buffer is converted into a byte slice and returned along with any encountered errors.
//
// Parameters:
// s vaulttypes.Vault: The Vault to be encrypted.
//
// Return:
// []byte: The encrypted byte slice.
// error: An error that occurred during the encryption process, or nil if no error occurred.
func (c *Cryptographer) Encrypt(s vaulttypes.Vault) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	encoded, err := vaulttypes.Serialise(s)
	if err != nil {
		return nil, err
	}

	context := Cryptographer{encoded}

	err = enc.Encode(context)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// Decrypt takes a byte slice as input and decrypts it using the Gob encoding.
// The decrypted data is expected to be in the form of a Cryptographer struct containing a byte slice representing a Vault.
// The function then decodes the Vault from the byte slice and returns it along with any encountered errors.
//
// Parameters:
// b []byte: The byte slice to be decrypted.
//
// Return:
// vaulttypes.Vault: The decrypted Vault.
// error: An error that occurred during the decryption process, or nil if no error occurred.
func (c *Cryptographer) Decrypt(b []byte) (vaulttypes.Vault, error) {
	var buff bytes.Buffer
	buff.Write(b)

	dec := gob.NewDecoder(&buff)

	var context Cryptographer
	err := dec.Decode(&context)
	if err != nil {
		return nil, err
	}
	return vaulttypes.Deserialise(context.Context)
}
