//go:build unit
// +build unit

package unit_test

import (
	"errors"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestProjectDeleteOptionsValidate(t *testing.T) {
	t.Parallel()

	validOptions := domain.ProjectDeleteOptions{
		Target: domain.ProjectIdentifier{Name: "sample-app"},
		Scope:  domain.DeleteScopeMetadata,
	}

	if err := validOptions.Validate(); err != nil {
		t.Fatalf("Validate() valid case returned error: %v", err)
	}

	cases := []struct {
		name    string
		options domain.ProjectDeleteOptions
		wantErr string
		isErr   error
	}{
		{
			name:    "missing target fails",
			options: domain.ProjectDeleteOptions{Scope: domain.DeleteScopeMetadata},
			wantErr: "delete target is required",
		},
		{
			name:    "invalid scope fails",
			options: domain.ProjectDeleteOptions{Target: domain.ProjectIdentifier{Name: "sample-app"}, Scope: domain.DeleteScope(-1)},
			isErr:   domain.ErrInvalidDeleteScope,
		},
		{
			name: "missing backup destination fails",
			options: domain.ProjectDeleteOptions{
				Target: domain.ProjectIdentifier{Name: "sample-app"},
				Scope:  domain.DeleteScopeMetadata,
				Backup: &domain.BackupRequest{},
			},
			wantErr: "backup destination must not be empty",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.options.Validate()
			if tc.isErr != nil {
				if !errors.Is(err, tc.isErr) {
					t.Fatalf("expected error %v, got %v", tc.isErr, err)
				}
				return
			}

			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("expected error %q, got %v", tc.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestProjectDeleteOptionsRequiresBackup(t *testing.T) {
	t.Parallel()

	opts := domain.ProjectDeleteOptions{}
	if opts.RequiresBackup() {
		t.Fatal("RequiresBackup should be false when request is nil")
	}

	opts.Backup = &domain.BackupRequest{Destination: "/tmp/sample.zip"}
	if !opts.RequiresBackup() {
		t.Fatal("RequiresBackup should be true when backup request present")
	}
}

func TestDeleteScopeValidationHelpers(t *testing.T) {
	t.Parallel()

	validScopes := []domain.DeleteScope{
		domain.DeleteScopeMetadata,
		domain.DeleteScopeKeepFiles,
		domain.DeleteScopeEnvOnly,
		domain.DeleteScopeAll,
	}

	for _, scope := range validScopes {
		if !scope.IsValid() {
			t.Fatalf("expected scope %v to be valid", scope)
		}
	}

	invalidScopes := []domain.DeleteScope{
		domain.DeleteScope(-1),
		domain.DeleteScope(42),
	}

	for _, scope := range invalidScopes {
		if scope.IsValid() {
			t.Fatalf("expected scope %v to be invalid", scope)
		}
	}

	if !domain.DeleteScopeAll.IncludesWorkspace() {
		t.Fatal("DeleteScopeAll should include workspace")
	}

	for _, scope := range []domain.DeleteScope{domain.DeleteScopeMetadata, domain.DeleteScopeKeepFiles, domain.DeleteScopeEnvOnly} {
		if scope.IncludesWorkspace() {
			t.Fatalf("scope %v should not include workspace", scope)
		}
	}

	for _, scope := range validScopes {
		if !scope.IncludesEnvironment() {
			t.Fatalf("scope %v should include environment cleanup", scope)
		}
	}

	metadataExpectations := []struct {
		scope domain.DeleteScope
		want  bool
	}{
		{domain.DeleteScopeMetadata, true},
		{domain.DeleteScopeKeepFiles, true},
		{domain.DeleteScopeEnvOnly, false},
		{domain.DeleteScopeAll, true},
	}

	for _, tc := range metadataExpectations {
		if got := tc.scope.IncludesMetadata(); got != tc.want {
			t.Fatalf("scope %v metadata expectation: want %t got %t", tc.scope, tc.want, got)
		}
	}
}

func TestProjectIdentifierAndStatusHelpers(t *testing.T) {
	t.Parallel()

	var id domain.ProjectIdentifier
	if !id.IsZero() {
		t.Fatal("zero identifier should report IsZero=true")
	}

	id = domain.ProjectIdentifier{Name: "sample", Path: "/tmp/sample", RegistryID: "registry-id"}
	if id.IsZero() {
		t.Fatal("non-zero identifier should report IsZero=false")
	}

	status := domain.ProjectStatus{}
	if !status.IsActionable() {
		t.Fatal("default project status should be actionable")
	}

	status = domain.ProjectStatus{Locked: true}
	if status.IsActionable() {
		t.Fatal("locked project should not be actionable")
	}

	status = domain.ProjectStatus{ReadOnly: true}
	if status.IsActionable() {
		t.Fatal("read-only project should not be actionable")
	}

	status = domain.ProjectStatus{Locked: true, ReadOnly: true}
	if status.IsActionable() {
		t.Fatal("locked and read-only project should not be actionable")
	}
}
