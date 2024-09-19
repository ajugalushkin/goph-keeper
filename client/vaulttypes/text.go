package vaulttypes

import (
	"fmt"
)

type Text struct {
	Data string
}

// Type возвращает тип хранимой информации
func (t Text) Type() VaultType {
	return vaultTypeText
}

// String функция отображения приватной информации
func (t Text) String() string {
	return fmt.Sprintf("TextData: %s", t.Data)
}
