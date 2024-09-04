package cmd

import (
	"bytes"
	"encoding/gob"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
)

// keepCmd represents the keep command
var keepCmd = &cobra.Command{
	Use:   "keep",
	Short: "Manage user private data",
}

func init() {
	rootCmd.AddCommand(keepCmd)
}

type Data struct {
	Context []byte
}

func encryptSecret(s vaulttypes.Vault) ([]byte, error) {
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

func decryptSecret(b []byte) (vaulttypes.Vault, error) {
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
