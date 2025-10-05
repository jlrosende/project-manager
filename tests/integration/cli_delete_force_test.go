//go:build integration
// +build integration

package integration_test

import (
	"context"
	"testing"

	deletecmd "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/delete"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_ForceSkipsConfirmation(t *testing.T) {
	svc := &fakeDeleteService{
		result: &domain.ProjectDeleteResult{
			Scope:            domain.DeleteScopeMetadata,
			ArtifactsRemoved: []domain.DeletionArtifact{{Type: domain.ArtifactRegistry}},
			Summary:          "forced",
		},
	}

	called := false
	confirm := func(_ context.Context, _ deletecmd.ConfirmationRequest) (bool, error) {
		called = true
		return true, nil
	}

	code, _, stderr, err := executeDeleteCommand(t, svc, confirm, "sample-app", "--force")
	if err != nil {
		t.Fatalf("expected success, got error: %v (stderr=%s)", err, stderr)
	}

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	if called {
		t.Fatal("confirmation should be skipped when --force is provided")
	}

	if svc.lastOpt.Force != true {
		t.Fatalf("expected Force flag true in options, got %+v", svc.lastOpt)
	}
}
