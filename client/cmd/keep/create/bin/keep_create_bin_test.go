package bin

import (
	"log/slog"
	"os"
	"testing"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
)

func TestNewCommand_NilLogger(t *testing.T) {
	t.Parallel()

	// Arrange
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var mockClient app.KeeperClient

	// Act
	cmd := NewCommand(log, mockClient)

	// Assert
	if cmd == nil {
		t.Error("expected NewCommand to return a non-nil command")
	}
}
