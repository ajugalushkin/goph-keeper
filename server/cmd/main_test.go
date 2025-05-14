package main

import (
	testCfg "github.com/ajugalushkin/goph-keeper/server/config"
	"github.com/ajugalushkin/goph-keeper/server/internal/app"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	err := os.Mkdir("test", 0777)
	if err != nil {
		return
	}
	newConfig := testCfg.Config{
		Env:     "dev",
		Storage: testCfg.Storage{Path: "postgresql://praktikum:pass@localhost:5432/goph_keeper?sslmode=disable"},
		GRPC: testCfg.GRPC{
			Address: "localhost:50051",
			Timeout: time.Hour,
		},
		Token: testCfg.Token{
			TTL:    time.Hour,
			Secret: "secret",
		},
		Minio: testCfg.Minio{
			Endpoint: "localhost:9000",
			Username: "testuser",
			Password: "testpassword",
			SSL:      false,
			Bucket:   "testbucket",
		},
	}
	yamlData, err := yaml.Marshal(&newConfig)
	if err != nil {
		return
	}
	filePath := "test/config.yaml"
	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return
	}
	os.Setenv("SERVER_CONFIG", filePath)

	exitcode := m.Run()

	os.RemoveAll("test")
	os.Clearenv()
	os.Exit(exitcode)
}

func TestSetupLogger_DevEnvironment(t *testing.T) {
	env := "dev"

	logger := setupLogger(env)

	if logger == nil {
		t.Fatalf("Expected a logger instance, but got nil")
	}

	handler := logger.Handler()
	if handler == nil {
		t.Fatalf("Expected a handler, but got nil")
	}

	_, ok := handler.(*slog.TextHandler)
	if !ok {
		t.Fatalf("Expected a TextHandler, but got %T", handler)
	}
}

func TestSetupLogger_Prod(t *testing.T) {
	const envProd = "prod"

	logger := setupLogger(envProd)

	if logger == nil {
		t.Fatal("logger is nil")
	}

	handler := logger.Handler()
	if handler == nil {
		t.Fatal("logger handler is nil")
	}

	_, ok := handler.(*slog.TextHandler)
	if !ok {
		t.Fatalf("unexpected handler type: %T", handler)
	}
}

func TestInitApp_NilLogger(t *testing.T) {
	cfg := testCfg.MustLoad()

	// Set up a nil logger.
	var log *slog.Logger

	// Expect a panic when initializing the application with a nil logger.
	assert.Panics(t, func() {
		initApp(log, cfg)
	}, "expected initApp to panic with a nil logger")
}

func TestInitApp_NilConfig(t *testing.T) {
	// Set up a logger for the test.
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Set up a nil configuration for the test.
	var cfg *testCfg.Config

	// Expect a panic when initializing the application with a nil configuration.
	assert.Panics(t, func() {
		initApp(log, cfg)
	}, "expected initApp to panic with a nil configuration")
}

func TestRun_NilApplication(t *testing.T) {
	// Create a nil application pointer.
	var application *app.App

	// Expect a panic when running the application with a nil pointer.
	assert.Panics(t, func() {
		run(application)
	}, "expected run to panic with a nil application")
}

func TestWaitSignal_NilApplication(t *testing.T) {
	// Create a nil logger instance.
	var log *slog.Logger

	// Create a nil application instance.
	var application *app.App

	// Create a channel to receive OS signals.
	stop := make(chan os.Signal, 1)

	// Simulate SIGINT signal.
	signal.Notify(stop, syscall.SIGINT)

	// Start the waitSignal function in a separate goroutine.
	go func() {
		waitSignal(log, application)
	}()

	// Send SIGINT signal to the waitSignal function.
	stop <- syscall.SIGINT

	// Wait for the waitSignal function to finish.
	time.Sleep(100 * time.Millisecond)

	// Verify that the waitSignal function did not panic.
	if r := recover(); r != nil {
		t.Fatalf("waitSignal function panicked: %v", r)
	}
}
