package root

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
)

func TestInit_TokenStorageWithFileBasedStorageWhenFileDoesNotExist(t *testing.T) {
	// Arrange
	tokenFilePath := "nonexistent_token.txt"
	defer os.Remove(tokenFilePath) // Clean up after the test

	// Act
	TokenStorage := token_cache.GetToken()

	// Assert
	assert.NotNil(t, TokenStorage, "Token storage should not be nil")
	assert.IsType(t, &token_cache.FileStorage{}, TokenStorage, "Token storage should be of type FileStorage")
}

func execute(args string) string {

	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs(strings.Split(args, " "))
	rootCmd.Execute()

	return actual.String()
}

func Test_ExecuteParticularCommandDynamically(t *testing.T) {
	//run all the child commands reside under root command
	for _, child := range rootCmd.Commands() {

		//check if args exist or not
		childArgs, exists := child.Annotations["args"]
		if exists {
			actual := execute(childArgs)
			expected, exists := child.Annotations["output"]
			if !exists {
				t.Errorf("Output doesn't exists in [%s] command", child.Use)
			}
			assert.Equal(t, actual, expected, "actual is not expected")
		}

		if !child.HasSubCommands() {
			continue
		}

		for _, grandChild := range child.Commands() {
			grandChildArgs, exists := grandChild.Annotations["args"]
			if exists {
				actual := execute(grandChildArgs)
				expected, exists := grandChild.Annotations["output"]
				if !exists {
					t.Errorf("Output doesn't exists in [%s] command", child.Use)
				}
				assert.Equal(t, actual, expected, "actual is not expected")
			}
		}

	}

}
