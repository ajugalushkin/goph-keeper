package register

import (
	"testing"

	"github.com/spf13/cobra"
)

// Command is added to authCmd successfully
func TestCommandRegAddedToAuthCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "root"}
	authCmd := &cobra.Command{Use: "auth"}
	rootCmd.AddCommand(authCmd)

	registerCmd := &cobra.Command{Use: "register"}
	authCmd.AddCommand(registerCmd)

	found := false
	for _, cmd := range authCmd.Commands() {
		if cmd.Use == "register" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("register command was not added to authCmd")
	}
}
