package login

import (
	"log/slog"
	"os"
	"strings"
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

func TestNewCommand_EmptyEmailAndValidPassword(t *testing.T) {
	// Arrange
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	expectedPassword := "testpassword"

	// Act
	cmd := NewCommand(mockLog, mockClient)
	cmd.SetArgs([]string{"-e", "", "-p", expectedPassword})
	err := cmd.Execute()

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "email is required")
	mockClient.AssertNotCalled(t, "Login", mock.Anything, mock.Anything, mock.Anything)
}
func TestNewCommand_EmptyEmailAndValidPassword2(t *testing.T) {
	// Arrange
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	expectedPassword := "testpassword"

	// Act
	cmd := NewCommand(mockLog, mockClient)
	cmd.SetArgs([]string{"-p", expectedPassword})
	err := cmd.Execute()

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "required flag(s) \"email\" not set")
	mockClient.AssertNotCalled(t, "Login", mock.Anything, mock.Anything, mock.Anything)
}

func TestNewCommand_ValidEmailAndEmptyPassword(t *testing.T) {
	// Arrange
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	expectedEmail := "test@example.com"
	expectedPassword := ""

	// Act
	cmd := NewCommand(mockLog, mockClient)
	cmd.SetArgs([]string{"-e", expectedEmail, "-p", expectedPassword})
	err := cmd.Execute()

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "password is required")
	mockClient.AssertExpectations(t)
}

func TestNewCommand_ValidEmailAndEmptyPassword2(t *testing.T) {
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
	mockClient.AssertExpectations(t)
}

func TestNewCommand_ValidEmailAndPasswordExceedingMaxLength(t *testing.T) {
	// Arrange
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	// Create a valid email with maximum length
	expectedEmail := strings.Repeat("a", 254) + "@example.com"
	// Create a valid password with maximum length
	expectedPassword := strings.Repeat("a", 255)

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

func TestNewCommand_ValidEmailAndPasswordWithSpecialCharacters(t *testing.T) {
	// Arrange
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	expectedEmail := "test.email+special@example.com"
	expectedPassword := "test$pecial_password"

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
func TestNewCommand_UppercaseEmailAndPassword(t *testing.T) {
	// Arrange
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	expectedEmail := "TEST@EXAMPLE.COM"
	expectedPassword := "TESTPASSWORD"

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

func TestNewCommand_NonASCIICharactersInEmailAndPassword(t *testing.T) {
	// Arrange
	mockLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockClient := mocks.NewAuthClient(t)

	expectedEmail := "tëst@exämple.cöm"
	expectedPassword := "tëstpässwörð"

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
