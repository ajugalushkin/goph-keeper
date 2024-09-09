package vaulttypes

// Bin произвольные бинарные данные
type Bin struct {
	Data []byte
}

// Type возвращает тип хранимой информации
func (b Bin) Type() VaultType {
	return vaultTypeBin
}

// String функция отображения приватной информации
func (b Bin) String() string {
	return "BINARY DATA"
}
