package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type BucketStorage struct {
	mc *minio.Client
}

func NewBucketStorage(storagePath string) (*BucketStorage, error) {
	const op = "storage.minio.NewBucketStorage"

	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Подключение к Minio с использованием имени пользователя и пароля
	client, err := minio.New(config.AppConfig.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AppConfig.MinioRootUser, config.AppConfig.MinioRootPassword, ""),
		Secure: config.AppConfig.MinioUseSSL,
	})
	if err != nil {
		return err
	}

	// Установка подключения Minio
	minioClient := &BucketStorage{
		mc: client,
	}

	// Проверка наличия бакета и его создание, если не существует
	exists, err := minioClient.mc.BucketExists(ctx, config.AppConfig.BucketName)
	if err != nil {
		return err
	}
	if !exists {
		err := minioClient.mc.MakeBucket(ctx, config.AppConfig.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return minioClient, nil
}

func (m *BucketStorage) InitMinio() error {

	return nil
}
