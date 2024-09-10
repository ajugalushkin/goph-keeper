package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// Command is successfully added to authCmd
func TestCommandAddedToAuthCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "root"}
	authCmd := &cobra.Command{Use: "auth"}
	loginCmd := &cobra.Command{Use: "login"}

	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)

	found := false
	for _, cmd := range authCmd.Commands() {
		if cmd.Use == "login" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("loginCmd was not added to authCmd")
	}
}

// Error occurs when marking email flag as required
func TestErrorMarkingEmailFlagRequired(t *testing.T) {
	loginCmd := &cobra.Command{Use: "login"}
	loginCmd.Flags().StringP("email", "e", "", "User Email")

	err := loginCmd.MarkFlagRequired("email")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Simulate an error scenario by trying to mark a non-existent flag as required
	err = loginCmd.MarkFlagRequired("nonexistent")
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}
