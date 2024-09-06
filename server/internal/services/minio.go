package services

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ajugalushkin/goph-keeper/server/config"
	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

type Minio struct {
	log *slog.Logger
	mc  *minio.Client
	cfg *config.Minio
}

func NewMinioService(
	log *slog.Logger,
	cfg config.Minio,
) *Minio {
	ctx := context.Background()

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Username, cfg.Password, ""),
		Secure: cfg.SSL,
	})
	if err != nil {
		log.Error("Minio New Error", slog.String("error", err.Error()))
		return nil
	}

	exists, err := client.BucketExists(ctx, cfg.Buket)
	if err != nil {
		log.Error("Minio New Error", slog.String("error", err.Error()))
		return nil
	}
	if !exists {
		err := client.MakeBucket(ctx, cfg.Buket, minio.MakeBucketOptions{})
		if err != nil {
			log.Error("Minio New Error", slog.String("error", err.Error()))
			return nil
		}
	}

	return &Minio{
		log: log,
		mc:  client,
		cfg: &cfg,
	}
}

func (m Minio) Create(file models.FileData) error {
	reader := bytes.NewReader(file.Data)

	opts := minio.PutObjectOptions{
		ContentType: http.DetectContentType(file.Data),
	}

	_, err := m.mc.PutObject(context.Background(), m.cfg.Buket, file.FileName, reader, int64(len(file.Data)), opts)
	if err != nil {
		return fmt.Errorf("minio create file error: %s", err.Error())
	}

	return nil
}

func (m Minio) Get() {
	//TODO implement me
	panic("implement me")
}
