package minio

import (
	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	mc *minio.Client
}

func NewMinioClient(storagePath string) (*MinioStorage, error) {
	const op = "storage.minio.NewMinioClient"

	//// Создание контекста с возможностью отмены операции
	//ctx := context.Background()
	//
	//cfg := config.GetInstance().Config
	//
	//// Подключение к Minio с использованием имени пользователя и пароля
	//client, err := minio.New(config.AppConfig.MinioEndpoint, &minio.Options{
	//	Creds:  credentials.NewStaticV4(config.AppConfig.MinioRootUser, config.AppConfig.MinioRootPassword, ""),
	//	Secure: config.AppConfig.MinioUseSSL,
	//})
	//if err != nil {
	//	return err
	//}
	//
	//// Установка подключения Minio
	//minioClient := &MinioStorage{
	//	mc: client,
	//}
	//
	//// Проверка наличия бакета и его создание, если не существует
	//exists, err := minioClient.mc.BucketExists(ctx, config.AppConfig.BucketName)
	//if err != nil {
	//	return err
	//}
	//if !exists {
	//	err := minioClient.mc.MakeBucket(ctx, config.AppConfig.BucketName, minio.MakeBucketOptions{})
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//return minioClient, nil
	return nil, nil
}

func (m *MinioStorage) InitMinio() error {

	return nil
}
