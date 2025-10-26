//go:build integration
// +build integration

package integration_test

import (
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_DryRunBackupReportsPath(t *testing.T) {
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			if !opts.DryRun {
				t.Fatal("expected dry-run true")
			}
			if opts.Backup == nil {
				t.Fatal("expected backup request even in dry-run")
			}
			return &domain.ProjectDeleteResult{
				Plan: &domain.ProjectDeletePlan{
					Scope: domain.DeleteScopeAll,
				},
				Summary:    "dry run with backup",
				BackupPath: opts.Backup.Destination,
			}, nil
		},
	}

	destination := "/backups/sample-preview.zip"
	code, stdout, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app", "--dry-run", "--all", "--backup", "--backup-destination", destination)
	if err != nil {
		t.Fatalf("dry-run backup execute: %v (stderr=%s)", err, stderr)
	}

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if !strings.Contains(stdout, destination) {
		t.Fatalf("expected output to reference backup destination %q; got %s", destination, stdout)
	}

	if svc.calls != 1 {
		t.Fatalf("expected service invocation, got %d", svc.calls)
	}
}
