package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/ajugalushkin/goph-keeper/server/config"
	"github.com/ajugalushkin/goph-keeper/server/internal/app"
	"github.com/ajugalushkin/goph-keeper/server/internal/app/grpc/mocks"
)

// Application starts successfully with valid configuration
func TestApplicationStartsSuccessfully(t *testing.T) {
	cfg := &config.Config{
		Env: "dev",
		GRPC: config.GRPC{
			Address: "localhost:50051",
		},
		Token: config.Token{
			TTL:    time.Hour,
			Secret: "secret",
		},
		Storage: config.Storage{
			Path: "/tmp/storage",
		},
		Minio: config.Minio{
			Endpoint: "localhost:9000",
			Username: "minio",
			Password: "minio123",
			SSL:      false,
			Bucket:   "test-bucket",
		},
	}

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	application := app.New(log, cfg)

	mockGRPCServer := mocks.NewGrpcServer(t)
	application.GRPCSrv = mockGRPCServer

	mockGRPCServer.On("MustRun").Return().Once()
	mockGRPCServer.On("Stop").Return().Once()

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		time.Sleep(1 * time.Second)
		stop <- syscall.SIGINT
	}()

	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()
	log.Info("application stopped")

	mockGRPCServer.AssertExpectations(t)
}

// Configuration file is missing or unreadable
func TestConfigurationFileMissingOrUnreadable(t *testing.T) {
	viper.Reset()

	// Simulate missing config file by setting an invalid path
	viper.SetConfigFile("/invalid/path/to/config.yaml")

	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(error).Error(), "Unable to unmarshal config file")
		}
	}()

	cfg := config.MustLoad()

	assert.Nil(t, cfg)
}
