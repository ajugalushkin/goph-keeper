package main

import (
	"log/slog"
	"testing"
)

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
