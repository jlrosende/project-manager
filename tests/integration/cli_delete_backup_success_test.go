//go:build integration
// +build integration

package integration_test

import (
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_BackupAllSuccess(t *testing.T) {
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			if opts.Backup == nil {
				t.Fatal("expected backup request to be populated")
			}
			if opts.Scope != domain.DeleteScopeAll {
				t.Fatalf("expected DeleteScopeAll, got %v", opts.Scope)
			}
			return &domain.ProjectDeleteResult{
				Scope:      opts.Scope,
				BackupPath: "/backups/sample-app.zip",
				ArtifactsRemoved: []domain.DeletionArtifact{
					{Type: domain.ArtifactBackup, Description: "archive"},
					{Type: domain.ArtifactWorkspace, Description: "workspace"},
				},
				Summary: "Backed up and removed",
			}, nil
		},
	}

	code, stdout, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app", "--all", "--backup", "--force")
	if err != nil {
		t.Fatalf("backup delete: %v (stderr=%s)", err, stderr)
	}

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if !strings.Contains(stdout, "/backups/sample-app.zip") {
		t.Fatalf("stdout missing backup path: %s", stdout)
	}

	if svc.calls != 1 || !svc.lastOpt.Force {
		t.Fatalf("expected force flag to propagate, got %+v", svc.lastOpt)
	}
}
