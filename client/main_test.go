package main

import (
	"testing"

	"github.com/ajugalushkin/goph-keeper/client/cmd"
)

// Execute function runs without errors
func TestExecuteRunsWithoutErrors(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked with error: %v", r)
		}
	}()

	cmd.Execute()
}
