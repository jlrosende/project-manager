//go:build unit
// +build unit

package unit_test

import (
	"strings"
	"testing"

	deletecmd "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/delete"
)

func TestCLIDeleteArgsRequiresTarget(t *testing.T) {
	cmd := deletecmd.Command()

	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Fatal("expected error when no target provided")
	}
}

func TestCLIDeleteArgsRejectsMultipleScopeFlags(t *testing.T) {
	tests := []struct {
		name  string
		flags map[string]string
	}{
		{
			name: "all_and_keep_files",
			flags: map[string]string{
				"all":        "true",
				"keep-files": "true",
			},
		},
		{
			name: "all_and_only_env",
			flags: map[string]string{
				"all":      "true",
				"only-env": "true",
			},
		},
		{
			name: "keep_files_and_only_env",
			flags: map[string]string{
				"keep-files": "true",
				"only-env":   "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := deletecmd.Command()

			for key, value := range tt.flags {
				if err := cmd.Flags().Set(key, value); err != nil {
					t.Fatalf("set flag %s: %v", key, err)
				}
			}

			err := cmd.Args(cmd, []string{"project"})
			if err == nil {
				t.Fatalf("expected error for flags %v", tt.flags)
			}

			msg := err.Error()
			for key := range tt.flags {
				if !strings.Contains(msg, key) {
					t.Fatalf("error message %q does not reference flag %s", msg, key)
				}
			}
		})
	}
}

func TestCLIDeleteArgsAllowsSingleScopeFlag(t *testing.T) {
	tests := []struct {
		name string
		flag string
	}{
		{name: "all", flag: "all"},
		{name: "keep_files", flag: "keep-files"},
		{name: "only_env", flag: "only-env"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := deletecmd.Command()

			if err := cmd.Flags().Set(tt.flag, "true"); err != nil {
				t.Fatalf("set flag %s: %v", tt.flag, err)
			}

			if err := cmd.Args(cmd, []string{"project"}); err != nil {
				t.Fatalf("unexpected error with flag %s: %v", tt.flag, err)
			}
		})
	}
}
