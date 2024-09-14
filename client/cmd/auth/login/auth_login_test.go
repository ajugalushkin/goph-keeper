package login

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/goph-keeper/client/internal/app/mocks"
)

func TestNewCommand_ValidEmailAndPassword(t *testing.T) {
	// Arrange
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	expectedEmail := "test@example.com"
	expectedPassword := "testpassword"

	// Set expectations for mock objects
	mockClient.On("Login", mock.Anything, expectedEmail, expectedPassword).Return("testToken", nil)

	// Act
	cmd := NewCommand(mockLog, mockClient)
	cmd.SetArgs([]string{"-e", expectedEmail, "-p", expectedPassword})
	err := cmd.Execute()

	// Assert
	require.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestNewCommand_MissingRequiredFlags(t *testing.T) {
	// Arrange
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	// Act
	cmd := NewCommand(mockLog, mockClient)
	cmd.SetArgs([]string{"-p", "testpassword"})
	err := cmd.Execute()

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "required flag(s) \"email\" not set")
	mockClient.AssertNotCalled(t, "Login", mock.Anything, mock.Anything, mock.Anything)
}

func TestNewCommand_ValidEmailAndMissingPassword(t *testing.T) {
	// Arrange
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	expectedEmail := "test@example.com"

	// Act
	cmd := NewCommand(mockLog, mockClient)
	cmd.SetArgs([]string{"-e", expectedEmail})
	err := cmd.Execute()

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "required flag(s) \"password\" not set")
	mockClient.AssertNotCalled(t, "Login", mock.Anything, mock.Anything, mock.Anything)
}

func TestNewCommand_MissingEmailFlag(t *testing.T) {
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	expectedPassword := "testpassword"

	cmd := NewCommand(mockLog, mockClient)
	cmd.SetArgs([]string{"-p", expectedPassword})
	err := cmd.Execute()

	require.Error(t, err)
	require.Contains(t, err.Error(), "required flag(s) \"email\" not set")
	mockClient.AssertNotCalled(t, "Login", mock.Anything, mock.Anything, mock.Anything)
}