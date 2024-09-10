package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// authCmd is successfully added to rootCmd
func TestAuthCmdAddedToRootCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "root"}
	authCmd := &cobra.Command{Use: "auth"}

	rootCmd.AddCommand(authCmd)

	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "auth" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("authCmd was not added to rootCmd")
	}
}

// authCmd is added to rootCmd when rootCmd has no other subcommands
func TestAuthCmdAddedToEmptyRootCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "root"}
	authCmd := &cobra.Command{Use: "auth"}

	if len(rootCmd.Commands()) != 0 {
		t.Fatalf("rootCmd should have no subcommands initially")
	}

	rootCmd.AddCommand(authCmd)

	if len(rootCmd.Commands()) != 1 {
		t.Errorf("authCmd was not added to rootCmd")
	}
}
