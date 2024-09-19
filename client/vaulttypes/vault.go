package vaulttypes

import (
	"encoding/json"
	"errors"
)

type VaultType string

const (
	vaultTypeCredentials VaultType = "creds"
	vaultTypeText        VaultType = "text"
	vaultTypeBin         VaultType = "bin"
	vaultTypeCard        VaultType = "card"
)

type Vault interface {
	Type() VaultType
	String() string
}

type container struct {
	Type VaultType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

func Serialise(vault Vault) ([]byte, error) {
	if vault == nil {
		return nil, errors.New("cannot encode nil vault")
	}
	data, err := json.Marshal(vault)
	if err != nil {
		return nil, err
	}
	return json.Marshal(container{
		Type: vault.Type(),
		Data: data,
	})
}

func Deserialise(data []byte) (Vault, error) {
	var c container
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	switch c.Type {
	case vaultTypeCredentials:
		var credentials Credentials
		if err := json.Unmarshal(c.Data, &credentials); err != nil {
			return nil, err
		}
		return credentials, nil
	case vaultTypeText:
		var text Text
		if err := json.Unmarshal(c.Data, &text); err != nil {
			return nil, err
		}
		return text, nil
	case vaultTypeBin:
		var bin Bin
		if err := json.Unmarshal(c.Data, &bin); err != nil {
			return nil, err
		}
		return bin, nil
	case vaultTypeCard:
		var card Card
		if err := json.Unmarshal(c.Data, &card); err != nil {
			return nil, err
		}
		return card, nil
	default:
		return nil, errors.New("unknown secret type")
	}
}
