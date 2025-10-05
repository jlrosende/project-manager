package unit_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestValidateProjectDefinition_SuccessCases(t *testing.T) {
	t.Run("absolute path remains", func(t *testing.T) {
		dir := t.TempDir()
		def := &domain.ProjectDefinition{Name: "ValidName", Path: dir}

		if err := domain.ValidateProjectDefinition(def); err != nil {
			t.Fatalf("validate project definition: %v", err)
		}

		if def.Path != dir {
			t.Fatalf("expected path %q, got %q", dir, def.Path)
		}
	})

	t.Run("relative path normalized", func(t *testing.T) {
		rel := filepath.Join("relative", "path")
		def := &domain.ProjectDefinition{Name: "Example", Path: rel}

		if err := domain.ValidateProjectDefinition(def); err != nil {
			t.Fatalf("validate project definition: %v", err)
		}

		want, err := filepath.Abs(rel)
		if err != nil {
			t.Fatalf("abs path: %v", err)
		}

		if def.Path != want {
			t.Fatalf("expected normalized path %q, got %q", want, def.Path)
		}
	})

	t.Run("here destination allowed", func(t *testing.T) {
		def := &domain.ProjectDefinition{Name: "Example", Here: true}
		if err := domain.ValidateProjectDefinition(def); err != nil {
			t.Fatalf("validate project definition: %v", err)
		}

		if def.Path != "" {
			t.Fatalf("expected empty path when using here, got %q", def.Path)
		}
	})
}

func TestValidateProjectDefinition_ValidationErrors(t *testing.T) {
	temp := t.TempDir()

	cases := []struct {
		name      string
		def       *domain.ProjectDefinition
		wantField map[string]string
	}{
		{
			name:      "missing name",
			def:       &domain.ProjectDefinition{Path: temp},
			wantField: map[string]string{"name": "must be provided"},
		},
		{
			name:      "invalid name",
			def:       &domain.ProjectDefinition{Name: "bad name", Path: temp},
			wantField: map[string]string{"name": "pattern"},
		},
		{
			name: "conflicting destination",
			def:  &domain.ProjectDefinition{Name: "Valid", Path: temp, Here: true},
			wantField: map[string]string{
				"destination": "cannot set both",
			},
		},
		{
			name: "missing destination",
			def:  &domain.ProjectDefinition{Name: "Valid"},
			wantField: map[string]string{
				"destination": "must provide either",
			},
		},
		{
			name: "blank path treated as missing",
			def:  &domain.ProjectDefinition{Name: "Valid", Path: "   "},
			wantField: map[string]string{
				"destination": "must provide either",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := domain.ValidateProjectDefinition(tc.def)
			if err == nil {
				t.Fatalf("expected error for %s", tc.name)
			}

			var verr domain.ProjectValidationErrors
			if !errors.As(err, &verr) {
				t.Fatalf("expected ProjectValidationErrors, got %T", err)
			}

			for field, contains := range tc.wantField {
				if !containsMessage(verr, field, contains) {
					t.Fatalf("expected field %s to contain %q in %v", field, contains, verr)
				}
			}
		})
	}
}

func containsMessage(errs domain.ProjectValidationErrors, field, fragment string) bool {
	for _, e := range errs {
		if e.Field == field && strings.Contains(e.Message, fragment) {
			return true
		}
	}

	return false
}

func TestValidateProjectDefinitionNil(t *testing.T) {
	err := domain.ValidateProjectDefinition(nil)
	if err == nil {
		t.Fatal("expected error when definition is nil")
	}

	if err.Error() != "project definition is nil" {
		t.Fatalf("unexpected error: %v", err)
	}
}
