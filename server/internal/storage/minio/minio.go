package minio

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ajugalushkin/goph-keeper/server/config"
	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

type MinioStorage struct {
	mc  *minio.Client
	cfg *config.Minio
}

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

	objectID := uuid.New().String()

	isExists, err := m.mc.BucketExists(ctx, m.cfg.Bucket)
	if err != nil || !isExists {
		err := m.mc.MakeBucket(ctx, m.cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			slog.Error("Minio New Error", slog.String("op", op), slog.String("error", err.Error()))
			return "", err
		}
	}

	opts := minio.PutObjectOptions{}

	_, err = m.mc.PutObject(ctx, m.cfg.Bucket, objectID, file.Data, file.Size, opts)
	if err != nil {
		return "", fmt.Errorf("error uploading file: %v", err)
	}

	return objectID, nil
}

func (m *MinioStorage) Get(ctx context.Context, objectID string) (*models.File, error) {
	const op = "minio.storage.Minio.Get"

	opts := minio.GetObjectOptions{}

	object, err := m.mc.GetObject(ctx, m.cfg.Bucket, objectID, opts)
	if err != nil {
		return nil, fmt.Errorf("op: %s, error uploading file: %v", op, err)
	}

	stat, err := object.Stat()
	if err != nil {
		return nil, fmt.Errorf("op: %s, error get stat file: %v", op, err)
	}

	return &models.File{
		Size: stat.Size,
		Data: object,
	}, nil
}
