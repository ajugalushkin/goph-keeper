package creds

import (
	"errors"
	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
	"github.com/ajugalushkin/goph-keeper/client/secret"
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
)

// Fails to retrieve 'name' flag from command-line arguments
func TestFailToRetrieveNameFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("login", "testLogin", "login flag")
	cmd.Flags().String("password", "testPassword", "password flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

// Fails to retrieve 'login' flag from command-line arguments
func TestFailToRetrieveLoginFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name flag")
	cmd.Flags().String("password", "testPassword", "password flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected an error for failing to retrieve 'login' flag, but no error was returned")
	}
}

// Fails to retrieve 'password' flag from command-line arguments
func TestFailToRetrievePasswordFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name flag")
	cmd.Flags().String("login", "testLogin", "login flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error for failing to retrieve 'password' flag, but got no error")
	}
}

func TestCreateCredentialsWithEmptyNameFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "name flag")
	cmd.Flags().String("login", "testLogin", "login flag")
	cmd.Flags().String("password", "testPassword", "password flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error for empty 'name' flag, but no error was returned")
	}
}
func TestEmptyLoginFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name flag")
	cmd.Flags().String("login", "", "login flag")
	cmd.Flags().String("password", "testPassword", "password flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error for empty 'login' flag, but no error was returned")
	}
}
func TestEmptyPasswordFlag(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	cmd := &cobra.Command{}
	cmd.Flags().String("name", "testName", "name flag")
	cmd.Flags().String("login", "testLogin", "login flag")
	cmd.Flags().String("password", "", "password flag")

	err := createCredentialsCmdRunE(cmd, []string{})
	if err == nil {
		t.Fatalf("Expected error for empty 'password' flag, but got nil")
	}
}

func TestCreateCredsSuccess(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	name := gofakeit.Name()
	login := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 8)

	cmd := &cobra.Command{}
	cmd.Flags().String("name", name, "name flag")
	cmd.Flags().String("login", login, "login flag")
	cmd.Flags().String("password", password, "password flag")

	mockClient := mocks.NewKeeperClient(t)

	credentials := vaulttypes.Credentials{
		Login:    login,
		Password: password,
	}

	content, err := secret.NewCryptographer().Encrypt(credentials)
	require.NoError(t, err)

	resp := &v1.CreateItemResponseV1{
		Name:    name,
		Version: "1",
	}

	mockClient.On("CreateItem", mock.Anything, &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	}).Return(resp, nil)
	initClient(mockClient)

	err = createCredentialsCmdRunE(cmd, []string{})
	require.NoError(t, err)
	assert.Equal(t, name, resp.Name)
}

func TestCreateCredsError(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	name := gofakeit.Name()
	login := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 8)

	cmd := &cobra.Command{}
	cmd.Flags().String("name", name, "name flag")
	cmd.Flags().String("login", login, "login flag")
	cmd.Flags().String("password", password, "password flag")

	mockClient := mocks.NewKeeperClient(t)

	credentials := vaulttypes.Credentials{
		Login:    login,
		Password: password,
	}

	content, err := secret.NewCryptographer().Encrypt(credentials)
	require.NoError(t, err)

	resp := &v1.CreateItemResponseV1{
		Name:    name,
		Version: "1",
	}

	mockClient.On("CreateItem", mock.Anything, &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	}).Return(resp, errors.New("expected error"))
	initClient(mockClient)

	err = createCredentialsCmdRunE(cmd, []string{})
	require.Error(t, err)
}

func TestCreateCredsError2(t *testing.T) {
	logger.InitLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), &config.Config{Env: "dev"})

	name := gofakeit.Name()
	login := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 8)

	cmd := &cobra.Command{}
	cmd.Flags().String("name", name, "name flag")
	cmd.Flags().String("login", login, "login flag")
	cmd.Flags().String("password", password, "password flag")

	mockClient := mocks.NewKeeperClient(t)

	credentials := vaulttypes.Credentials{
		Login:    login,
		Password: password,
	}

	content, err := secret.NewCryptographer().Encrypt(credentials)
	require.NoError(t, err)

	mockClient.On("CreateItem", mock.Anything, &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	}).Return(nil, nil)
	initClient(mockClient)

	err = createCredentialsCmdRunE(cmd, []string{})
	require.Error(t, err)
}
