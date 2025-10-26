//go:build contract
// +build contract

package contract_test

import (
	"bytes"
	"strings"
	"testing"

	cli "github.com/jlrosende/project-manager/internal/adapters/handlers/cli"
)

func TestCLIDelete_HelpShowsExpectedFlags(t *testing.T) {
	cmd := cli.Root()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"delete", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("pm delete --help: %v", err)
	}

	output := out.String()
	expectedPhrases := []string{
		"pm delete <target>",
		"--all",
		"--keep-files",
		"--only-env",
		"--dry-run",
		"--force",
		"--backup",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Fatalf("help output missing %q; got: %s", phrase, output)
		}
	}
}
