package list

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

func TestKeepListRunE_EmptyList(t *testing.T) {
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
