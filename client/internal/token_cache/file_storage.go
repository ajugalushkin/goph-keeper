package token_cache

import (
	"io"
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
			if storage == nil {
				storage = newFileStorage("token_cache.txt")
			}
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

	defer file.Close()

	_, err = file.WriteString(accessToken)
	return err
}

// Load читает токен из файла
func (s *FileStorage) Load() (string, error) {
	file, err := os.Open(s.Path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if len(b) == 0 {
		return "", io.EOF
	}

	return string(b), err
}
