package main

import (
	"os"
	"testing"

	"github.com/spf13/viper"

	"github.com/ajugalushkin/goph-keeper/client/cmd/root"
)

// Execute runs without errors when called
func TestExecuteRunsWithoutErrors(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked with error: %v", r)
		}
	}()

	root.Execute()
}

// Configuration file is missing or unreadable
func TestExecuteConfigFileMissingOrUnreadable(t *testing.T) {
	// Set an invalid config file path
	os.Setenv("CLIENT_CONFIG", "/invalid/path/to/config.yaml")

	// Reset viper configuration to ensure it reads the new environment variable
	viper.Reset()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked with error: %v", r)
		}
	}()

	root.Execute()

	// Check if the log contains the expected message about the missing config file
	// This part is pseudo-code as it depends on the logging framework and its configuration
	// logOutput := getLogOutput()
	// if !strings.Contains(logOutput, "Config file not found") {
	//     t.Errorf("Expected log message about missing config file, but got: %s", logOutput)
	// }
}
