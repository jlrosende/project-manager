//go:build unit
// +build unit

package unit_test

import (
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
		name string
		args []string
	}{
		{name: "all_and_keep_files", args: []string{"project", "--all", "--keep-files"}},
		{name: "all_and_only_env", args: []string{"project", "--all", "--only-env"}},
		{name: "keep_files_and_only_env", args: []string{"project", "--keep-files", "--only-env"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &stubDeleteService{}
			code, _, _, err := deletecmd.ExecuteForTesting(deletecmd.Options{Service: service}, tt.args)
			if err == nil {
				t.Fatalf("expected error for args %v", tt.args)
			}

			if code == 0 {
				t.Fatalf("expected non-zero exit code for args %v", tt.args)
			}
		})
	}
}

func TestCLIDeleteArgsAllowsSingleScopeFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "all", args: []string{"project", "--all"}},
		{name: "keep_files", args: []string{"project", "--keep-files"}},
		{name: "only_env", args: []string{"project", "--only-env"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &stubDeleteService{}
			code, _, _, err := deletecmd.ExecuteForTesting(deletecmd.Options{Service: service}, tt.args)
			if err != nil {
				t.Fatalf("unexpected error with args %v: %v", tt.args, err)
			}

			if code != 0 {
				t.Fatalf("expected zero exit code, got %d", code)
			}
		})
	}
}
