package token

import (
	"io"
	"log/slog"
	"os"
)

// FileStorage файловое хранилище для токена
type FileStorage struct {
	Path string
}

// NewFileStorage создает новое файловое хранилище для токена
func NewFileStorage(path string) *FileStorage {
	return &FileStorage{
		Path: path,
	}
}

// Save записывает токен в файл
func (s *FileStorage) Save(accessToken string) error {
	file, err := os.Create(s.Path)
	if err != nil {
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			slog.Error("Error closing file: ", slog.String("error", err.Error()))
		}
	}()

	_, err = file.WriteString(accessToken)
	return err
}

// Load читает токен из файла
func (s *FileStorage) Load() (string, error) {
	file, err := os.Open(s.Path)
	if err != nil {
		return "", nil
	}
	defer func() {
		if err = file.Close(); err != nil {
			slog.Error("Error closing file: ", slog.String("error", err.Error()))
		}
	}()

	b, err := io.ReadAll(file)
	return string(b), err
}
