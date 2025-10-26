//go:build integration
// +build integration

package integration_test

import (
	"context"
	"strings"
	"testing"

	deletecmd "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/delete"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_NameSuccess(t *testing.T) {
	svc := &fakeDeleteService{
		result: &domain.ProjectDeleteResult{
			Scope: domain.DeleteScopeMetadata,
			ArtifactsRemoved: []domain.DeletionArtifact{
				{Type: domain.ArtifactRegistry, Description: "registry entry"},
				{Type: domain.ArtifactEnvVars, Description: "env vars"},
			},
			Summary: "Deleted project metadata",
		},
	}

	var confirmed bool
	confirm := func(_ context.Context, req deletecmd.ConfirmationRequest) (bool, error) {
		if req.Target.Name != "sample-app" {
			t.Fatalf("unexpected confirmation target: %+v", req.Target)
		}
		confirmed = true
		return true, nil
	}

	code, stdout, stderr, err := executeDeleteCommand(t, svc, confirm, "sample-app")
	if err != nil {
		t.Fatalf("execute delete: %v (stderr=%s)", err, stderr)
	}

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	if !confirmed {
		t.Fatal("expected confirmation prompt to run")
	}

	if svc.calls != 1 {
		t.Fatalf("expected DeleteProject invoked once, got %d", svc.calls)
	}

	if svc.lastOpt.Target.Name != "sample-app" {
		t.Fatalf("expected service to receive target name, got %+v", svc.lastOpt.Target)
	}

	if !strings.Contains(stdout, "registry entry") {
		t.Fatalf("stdout missing artifact summary: %s", stdout)
	}
}
