package auth

import (
	"reflect"
	"sort"
	"testing"

	"github.com/spf13/cobra"
)

func TestAuthCmdAddedToRootCmdWithSubcommands(t *testing.T) {
	rootCmd := &cobra.Command{Use: "root"}
	authCmd := NewCommand()

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

	expectedSubcommands := []string{"login", "register"}
	actualSubcommands := make([]string, 0, len(authCmd.Commands()))
	for _, c := range authCmd.Commands() {
		actualSubcommands = append(actualSubcommands, c.Use)
	}

	sort.Strings(expectedSubcommands)
	sort.Strings(actualSubcommands)

	if !reflect.DeepEqual(expectedSubcommands, actualSubcommands) {
		t.Errorf("authCmd subcommands do not match expected: %v, got: %v", expectedSubcommands, actualSubcommands)
	}
}
