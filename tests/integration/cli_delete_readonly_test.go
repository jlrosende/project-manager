//go:build integration
// +build integration

package integration_test

import (
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_ReadOnlyWorkspaceRejected(t *testing.T) {
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			return nil, domain.ErrWorkspaceReadOnly
		},
	}

	code, _, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app", "--all")
	if err == nil {
		t.Fatalf("expected error, got nil (code=%d)", code)
	}

	if code == 0 {
		t.Fatalf("expected non-zero exit for read-only, got %d", code)
	}

	if stderr == "" {
		t.Fatal("expected stderr message for read-only rejection")
	}
}

func TestCLIDelete_LockedProjectRejected(t *testing.T) {
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			return nil, domain.ErrProjectLocked
		},
	}

	code, _, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app")
	if err == nil {
		t.Fatalf("expected error, got nil (code=%d)", code)
	}

	if code == 0 {
		t.Fatalf("expected non-zero exit for locked project, got %d", code)
	}

	if stderr == "" {
		t.Fatal("expected stderr message for locked project")
	}
}
