package secret

import (
	"bytes"
	"encoding/gob"

	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
)

// Data is a struct used for encoding and decoding Vault data using Gob encoding.
type Data struct {
	Context []byte
}

// EncryptSecret takes a Vault as input and encrypts it using the Gob encoding.
// The function first encodes the Vault into a byte slice using the EncodeVault function from the vaulttypes package.
// Then, it creates a new Data struct containing the encoded byte slice.
// The Data struct is encoded using the Gob encoding into a bytes.Buffer.
// Finally, the bytes.Buffer is converted into a byte slice and returned along with any encountered errors.
//
// Parameters:
// s vaulttypes.Vault: The Vault to be encrypted.
//
// Return:
// []byte: The encrypted byte slice.
// error: An error that occurred during the encryption process, or nil if no error occurred.
func EncryptSecret(s vaulttypes.Vault) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	encoded, err := vaulttypes.EncodeVault(s)
	if err != nil {
		return nil, err
	}

	data := Data{encoded}

	err = enc.Encode(data)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// DecryptSecret takes a byte slice as input and decrypts it using the Gob encoding.
// The decrypted data is expected to be in the form of a Data struct containing a byte slice representing a Vault.
// The function then decodes the Vault from the byte slice and returns it along with any encountered errors.
//
// Parameters:
// b []byte: The byte slice to be decrypted.
//
// Return:
// vaulttypes.Vault: The decrypted Vault.
// error: An error that occurred during the decryption process, or nil if no error occurred.
func DecryptSecret(b []byte) (vaulttypes.Vault, error) {
	var buff bytes.Buffer
	buff.Write(b)

	dec := gob.NewDecoder(&buff)

	var data Data
	err := dec.Decode(&data)
	if err != nil {
		return nil, err
	}
	return vaulttypes.DecodeVault(data.Context)
}
