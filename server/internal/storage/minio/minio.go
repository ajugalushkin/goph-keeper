package minio

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/server/config"
	"log/slog"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
)

type FileStorage struct {
	mc  *minio.Client
	cfg *config.Minio
}

// NewMinioStorage initializes a new instance of FileStorage using Minio as the storage backend.
// It takes a Minio configuration as input and returns a pointer to the FileStorage instance and an error if any.
//
// The function performs the following steps:
// 1. Creates a new Minio client using the provided Minio configuration.
// 2. Checks if the specified bucket exists in Minio. If not, it creates the bucket.
// 3. Returns a pointer to a new FileStorage instance with the Minio client and configuration.
//
// Parameters:
// - cfg: A Minio configuration object containing the necessary details for connecting to the Minio server.
//
// Return:
// - A pointer to a new FileStorage instance or an error if any.
func NewMinioStorage(
	cfg config.Minio,
) (*FileStorage, error) {
	ctx := context.Background()

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Username, cfg.Password, ""),
		Secure: cfg.SSL,
	})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &FileStorage{
		mc:  client,
		cfg: &cfg,
	}, nil
}

// Create uploads a file to the Minio storage backend and returns the unique object ID.
//
// The function performs the following steps:
// 1. Generates a new UUID for the object ID.
// 2. Checks if the specified bucket exists in Minio. If not, it creates the bucket.
// 3. Uploads the file data to Minio using the generated object ID.
// 4. Returns the object ID and any encountered error.
//
// Parameters:
// - ctx: A context object that carries deadlines, cancellation signals, and other request-scoped values across API boundaries.
// - file: A pointer to a File model object containing the file data and size.
//
// Return:
// - A string representing the unique object ID of the uploaded file.
// - An error if any occurred during the file upload process.
func (m *FileStorage) Create(
	ctx context.Context,
	file *models.File,
) (string, error) {
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

// Get retrieves a file from the Minio storage backend using the provided object ID.
//
// The function performs the following steps:
// 1. Constructs GetObjectOptions for the Minio client.
// 2. Calls the GetObject method of the Minio client to retrieve the file data.
// 3. Calls the Stat method of the retrieved object to get its size.
// 4. Returns a pointer to a File model object containing the file data and size, or an error if any occurred.
//
// Parameters:
// - ctx: A context object that carries deadlines, cancellation signals, and other request-scoped values across API boundaries.
// - objectID: A string representing the unique identifier of the file to be retrieved.
//
// Return:
// - A pointer to a File model object containing the file data and size.
// - An error if any occurred during the retrieval process.
func (m *FileStorage) Get(ctx context.Context, objectID string) (*models.File, error) {
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

// Delete removes a file from the Minio storage backend using the provided object ID.
//
// The function takes a context object and an object ID as input parameters. It uses the Minio client
// to remove the specified object from the configured bucket. If an error occurs during the deletion
// process, it returns an error. Otherwise, it returns nil.
//
// Parameters:
// - ctx: A context object that carries deadlines, cancellation signals, and other request-scoped values across API boundaries.
// - objectID: A string representing the unique identifier of the file to be deleted.
//
// Return:
// - An error if any occurred during the deletion process. If the deletion is successful, it returns nil.
func (m *FileStorage) Delete(ctx context.Context, objectID string) error {
	const op = "minio.storage.Minio.Delete"

	err := m.mc.RemoveObject(ctx, m.cfg.Bucket, objectID, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("op: %s, error deleting file: %v", op, err)
	}

	return nil
}
