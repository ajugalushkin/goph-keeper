package token_cache

import (
	"io"
	"log/slog"
	"os"
	"sync"
)

var (
	storage Storage
	once    sync.Once
)

func InitTokenStorage(path string) {
	storage = newFileStorage(path)
}

func GetToken() Storage {
	once.Do(
		func() {
			storage = newFileStorage("token_cache.txt")
		})

	return storage
}

// FileStorage файловое хранилище для токена
type FileStorage struct {
	Path string
}

// NewFileStorage создает новое файловое хранилище для токена
func newFileStorage(path string) *FileStorage {
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
