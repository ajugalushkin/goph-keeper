package minio

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ajugalushkin/goph-keeper/server/config"
	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

type MinioStorage struct {
	mc  *minio.Client
	cfg *config.Minio
}

const bucketNameTemplate = "bucket-%d"

func NewMinioStorage(
	cfg config.Minio,
) (*MinioStorage, error) {
	const op = "storage.minio.NewMinioStorage"
	ctx := context.Background()

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Username, cfg.Password, ""),
		Secure: cfg.SSL,
	})
	if err != nil {
		slog.Error("Minio New Error", slog.String("error", err.Error()))
		return nil, err
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		slog.Error("Minio New Error", slog.String("error", err.Error()))
		return nil, err
	}
	if !exists {
		err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			slog.Error("Minio New Error", slog.String("error", err.Error()))
			return nil, err
		}
	}

	return &MinioStorage{
		mc:  client,
		cfg: &cfg,
	}, nil
}

func (m *MinioStorage) Create(ctx context.Context, file *models.File) (string, error) {
	const op = "minio.storage.Minio.Create"

	bucketName := fmt.Sprintf(bucketNameTemplate, file.UserID)
	isExists, err := m.mc.BucketExists(ctx, bucketName)
	if err != nil || !isExists {
		err := m.mc.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			slog.Error("Minio New Error", slog.String("op", op), slog.String("error", err.Error()))
			return "", err
		}
	}

	opts := minio.PutObjectOptions{
		ContentType: file.Type,
	}

	fileInfo, err := m.mc.PutObject(ctx, bucketName, file.Name, file.Data, file.Size, opts)
	if err != nil {
		return "", fmt.Errorf("error uploading file: %v", err)
	}

	return fileInfo.VersionID, nil
}

func (m *MinioStorage) Get(ctx context.Context, userID int64, fileName string) (*minio.Object, error) {
	opts := minio.GetObjectOptions{}

	bucketName := fmt.Sprintf(bucketNameTemplate, userID)
	object, err := m.mc.GetObject(ctx, bucketName, fileName, opts)
	if err != nil {
		return nil, fmt.Errorf("error uploading file: %v", err)
	}

	return object, nil
}
