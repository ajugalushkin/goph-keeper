package vaulttypes

import (
	"fmt"
)

// Credentials пара логин пароль
type Credentials struct {
	Login    string
	Password string
}

// Type возвращает тип хранимой информации
func (c Credentials) Type() VaultType {
	return vaultTypeCredentials
}

// String функция отображения приватной информации
func (c Credentials) String() string {
	return fmt.Sprintf("Login: %s, Password: %s", c.Login, c.Password)
}
