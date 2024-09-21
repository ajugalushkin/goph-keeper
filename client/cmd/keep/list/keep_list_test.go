package list

import (
	"context"
	"errors"
	"github.com/ajugalushkin/goph-keeper/mocks"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// Command creation with correct use and short description
func TestNewCommandCreation(t *testing.T) {
	initClient(nil)
	cmd := NewCommand()

	if cmd.Use != "list" {
		t.Errorf("expected command use to be 'list', got %s", cmd.Use)
	}

	if cmd.Short != "List secrets" {
		t.Errorf("expected command short description to be 'List secrets', got %s", cmd.Short)
	}

	if cmd.RunE == nil {
		t.Error("expected command RunE to be set, got nil")
	}
}

func TestKeepListRunE_EmptyList(t *testing.T) {
	initClient(nil)
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	mockClient.On(
		"ListItems",
		mock.Anything,
		&v1.ListItemsRequestV1{},
	).Return(&v1.ListItemsResponseV1{Secrets: nil}, nil)

	err := keepListRunE(nil, nil)
	assert.NoError(t, err)
}

func TestKeepListRunE_Error(t *testing.T) {
	initClient(nil)
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	expectedErr := errors.New("failed to list secrets")
	mockClient.On(
		"ListItems",
		mock.Anything,
		&v1.ListItemsRequestV1{},
	).Return(nil, expectedErr)

	err := keepListRunE(nil, nil)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestNewCommand_NoItemsFound(t *testing.T) {
	initClient(nil)
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	expectedErr := errors.New("no items found")
	mockClient.On(
		"ListItems",
		mock.Anything,
		&v1.ListItemsRequestV1{},
	).Return(nil, expectedErr)

	cmd := NewCommand()
	cmd.SetArgs([]string{"list"})
	err := cmd.Execute()
	assert.EqualError(t, err, expectedErr.Error())
}

// Successful execution of keepListRunE function
func TestKeepListRunEErrorHandling(t *testing.T) {
	initClient(nil)
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	// Create a new command
	cmd := NewCommand()

	// Mock the client
	mockClient := mocks.NewKeeperClient(t)
	initClient(mockClient)

	// Mock an error response from ListItems
	mockClient.On("ListItems", mock.Anything, mock.Anything).Return(nil, errors.New("mocked error"))

	// Call the function under test
	err := keepListRunE(cmd, []string{})

	// Check if the error is handled correctly
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

// Logging of errors during secret listing
func TestKeepListRunErrorLogging(t *testing.T) {
	initClient(nil)
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		&config.Config{Env: "dev"})

	// Create a fake cobra.Command
	fakeCmd := &cobra.Command{}

	mockError := errors.New("simulated error")
	mockClient := mocks.NewKeeperClient(t)
	mockClient.On(
		"ListItems",
		context.Background(),
		&v1.ListItemsRequestV1{},
	).Return(&v1.ListItemsResponseV1{}, mockError)
	initClient(mockClient)

	// Call the function under test
	err := keepListRunE(fakeCmd, []string{})

	// Check if the error is logged
	if !strings.Contains(err.Error(), mockError.Error()) {
		t.Errorf("expected error to be logged, got: %s", err.Error())
	}

	// Check if the error is returned
	if err == nil || err.Error() != mockError.Error() {
		t.Errorf("expected error to be returned, got: %v", err)
	}
}
