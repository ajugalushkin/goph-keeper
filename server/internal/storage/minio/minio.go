package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ajugalushkin/goph-keeper/server/config"
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

func (m *MinioStorage) Create(ctx context.Context, bucketName string, chunkNumber int, chunkData []byte) error {
	const op = "minio.storage.Minio.Create"

	isExists, err := m.mc.BucketExists(ctx, bucketName)
	if err != nil || !isExists {
		err := m.mc.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			slog.Error("Minio New Error", slog.String("op", op), slog.String("error", err.Error()))
			return err
		}
	}

	objectID := "Chunk-" + strconv.Itoa(chunkNumber)

	reader := bytes.NewReader(chunkData)

	contentType := http.DetectContentType(chunkData)

	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err = m.mc.PutObject(ctx, bucketName, objectID, reader, int64(len(chunkData)), opts)
	if err != nil {
		return fmt.Errorf("error uploading file: %v", err)
	}

	return nil
}

func (m *MinioStorage) Get(ctx context.Context, fileID string) (bytes.Buffer, error) {
	opts := minio.GetObjectOptions{}

	object, err := m.mc.GetObject(ctx, m.cfg.Bucket, fileID, opts)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("error uploading file: %v", err)
	}
	defer object.Close()

	fileData := bytes.Buffer{}
	buff := make([]byte, 1024)
	for {
		n, err := object.Read(buff)
		if err != nil && err != io.EOF {
			break
		}

		_, err = fileData.Write(buff[:n])
		if err != nil {
			return bytes.Buffer{}, fmt.Errorf("error uploading file: %v", err)
		}
	}
	return fileData, nil
}
