package root

import (
	"bytes"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
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
	err := rootCmd.Execute()
	if err != nil {
		return ""
	}

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

/*func TestInitConfig_WithNonExistentConfigFile(t *testing.T) {
	// Arrange
	nonExistentFilePath := "nonexistent_config.yaml"
	os.Setenv("CLIENT_CONFIG", nonExistentFilePath)
	defer os.Unsetenv("CLIENT_CONFIG")

	// Act
	initConfig()

	// Assert
	assert.Equal(t, nonExistentFilePath, viper.ConfigFileUsed(), "Config file path should match the CLIENT_CONFIG environment variable")
}*/
