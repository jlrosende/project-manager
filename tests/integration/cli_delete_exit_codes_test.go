//go:build integration
// +build integration

package integration_test

import (
	"errors"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDeleteExitCodes_SuccessAndFailure(t *testing.T) {
	successSvc := &fakeDeleteService{
		result: &domain.ProjectDeleteResult{Summary: "ok"},
	}

	successCode, _, _, successErr := executeDeleteCommand(t, successSvc, nil, "sample-app")
	if successErr != nil {
		t.Fatalf("unexpected error on success case: %v", successErr)
	}

	if successCode != 0 {
		t.Fatalf("expected success exit code 0, got %d", successCode)
	}

	dryRunSvc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			return &domain.ProjectDeleteResult{Plan: &domain.ProjectDeletePlan{Scope: domain.DeleteScopeMetadata}}, nil
		},
	}

	dryCode, _, _, dryErr := executeDeleteCommand(t, dryRunSvc, nil, "sample-app", "--dry-run")
	if dryErr != nil {
		t.Fatalf("unexpected error on dry run: %v", dryErr)
	}

	if dryCode != 0 {
		t.Fatalf("expected dry-run exit code 0, got %d", dryCode)
	}

	expectedErr := errors.New("boom")
	failureSvc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			return nil, expectedErr
		},
	}

	failCode, _, _, failErr := executeDeleteCommand(t, failureSvc, nil, "sample-app")
	if failErr == nil {
		t.Fatalf("expected failure error, got nil (code=%d)", failCode)
	}

	if failCode == 0 {
		t.Fatalf("expected non-zero exit on failure, got %d", failCode)
	}

	if !errors.Is(failErr, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, failErr)
	}
}
