//go:build unit
// +build unit

package unit_test

import (
	"context"
	"errors"
	"testing"

	deletecmd "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/delete"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

type stubDeleteService struct {
	result *domain.ProjectDeleteResult
	err    error
	opts   []domain.ProjectDeleteOptions
}

func (s *stubDeleteService) DeleteProject(ctx context.Context, opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
	s.opts = append(s.opts, opts)
	return s.result, s.err
}

func TestCLIDeleteExitCodeSuccess(t *testing.T) {
	service := &stubDeleteService{
		result: &domain.ProjectDeleteResult{
			Scope:            domain.DeleteScopeMetadata,
			ArtifactsRemoved: []domain.DeletionArtifact{{Type: domain.ArtifactRegistry}},
		},
	}

	code, stdout, stderr, err := deletecmd.ExecuteForTesting(deletecmd.Options{
		Service: service,
		Confirm: func(context.Context, deletecmd.ConfirmationRequest) (bool, error) { return true, nil },
	}, []string{"sample-app"})

	if err != nil {
		t.Fatalf("execute delete: %v (stderr=%s)", err, stderr)
	}

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stdout=%s, stderr=%s)", code, stdout, stderr)
	}

	if len(service.opts) != 1 {
		t.Fatalf("expected DeleteProject invoked once, got %d", len(service.opts))
	}
}

func TestCLIDeleteExitCodeDryRun(t *testing.T) {
	service := &stubDeleteService{
		result: &domain.ProjectDeleteResult{
			Scope:   domain.DeleteScopeAll,
			Plan:    &domain.ProjectDeletePlan{Scope: domain.DeleteScopeAll},
			Summary: "dry run",
		},
	}

	code, stdout, stderr, err := deletecmd.ExecuteForTesting(deletecmd.Options{
		Service: service,
		Confirm: func(context.Context, deletecmd.ConfirmationRequest) (bool, error) { return true, nil },
	}, []string{"sample-app", "--dry-run", "--all"})

	if err != nil {
		t.Fatalf("execute dry-run: %v (%s)", err, stderr)
	}

	if code != 0 {
		t.Fatalf("dry-run should exit 0, got %d", code)
	}

	if len(service.opts) != 1 || !service.opts[0].DryRun {
		t.Fatalf("expected DryRun option true, got %+v", service.opts)
	}

	if len(stdout) == 0 {
		t.Fatal("expected dry-run output to be written")
	}
}

func TestCLIDeleteExitCodeFailure(t *testing.T) {
	executeErr := errors.New("service failed")
	service := &stubDeleteService{err: executeErr}

	code, stdout, stderr, err := deletecmd.ExecuteForTesting(deletecmd.Options{
		Service: service,
		Confirm: func(context.Context, deletecmd.ConfirmationRequest) (bool, error) { return true, nil },
	}, []string{"missing-app"})

	if err == nil {
		t.Fatalf("expected error, got nil (code=%d, stdout=%s, stderr=%s)", code, stdout, stderr)
	}

	if !errors.Is(err, executeErr) {
		t.Fatalf("expected wrapped service error, got %v", err)
	}

	if code == 0 {
		t.Fatalf("expected non-zero exit code on failure, got %d", code)
	}

	if len(stderr) == 0 {
		t.Fatal("expected stderr output on failure")
	}
}
