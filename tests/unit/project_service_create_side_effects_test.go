//go:build unit
// +build unit

package unit_test

import (
	"errors"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/mocks"
	"go.uber.org/mock/gomock"
)

type stubFilesystem struct {
	removed []string
}

func (stubFilesystem) EnsureDir(string, fs.FileMode) error         { return nil }
func (stubFilesystem) IsDirEmpty(string) (bool, error)             { return true, nil }
func (stubFilesystem) Rename(string, string) error                 { return nil }
func (stubFilesystem) WriteFile(string, []byte, fs.FileMode) error { return nil }

func (stubFilesystem) Join(elem ...string) string      { return filepath.Join(elem...) }
func (stubFilesystem) IsAbs(path string) bool          { return filepath.IsAbs(path) }
func (stubFilesystem) Abs(path string) (string, error) { return filepath.Abs(path) }
func (stubFilesystem) ExpandHome(path string) string   { return path }
func (s *stubFilesystem) Remove(path string) error     { s.removed = append(s.removed, path); return nil }
func (stubFilesystem) UserHomeDir() (string, error)    { return "/home/test", nil }

func TestProjectService_Create_CleansProjectOnEnvSaveFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	projRepo := mocks.NewMockProjectRepository(ctrl)
	envRepo := mocks.NewMockEnvVarsRepository(ctrl)
	gitRepo := mocks.NewMockGitRepository(ctrl)

	fs := &stubFilesystem{}

	proj := &domain.Project{Name: "svc", Path: "/tmp/svc", EnvVarsFile: ".env"}
	projRepo.EXPECT().Create("svc", "/tmp/svc", "", "sh", ".env", gomock.Any(), gomock.Any()).Return(proj, nil)
	envRepo.EXPECT().Save(filepath.Join(proj.Path, ".env"), gomock.Any()).Return(errors.New("boom"))

	svc := services.NewProjectService(projRepo, envRepo, gitRepo, fs, nil)

	_, err := svc.Create("svc", "/tmp/svc", "", "sh", ".env", domain.EnvVars{}, nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	expected := filepath.Join(proj.Path, ".project.hcl")
	if len(fs.removed) != 1 || fs.removed[0] != expected {
		t.Fatalf("expected remove of %s, got %#v", expected, fs.removed)
	}
}

func TestProjectService_Create_CleansArtifactsOnIncludeIfFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	projRepo := mocks.NewMockProjectRepository(ctrl)
	envRepo := mocks.NewMockEnvVarsRepository(ctrl)
	gitRepo := mocks.NewMockGitRepository(ctrl)

	fs := &stubFilesystem{}

	proj := &domain.Project{Name: "svc", Path: "/tmp/svc", EnvVarsFile: ".env"}
	projRepo.EXPECT().Create("svc", "/tmp/svc", "", "sh", ".env", gomock.Any(), gomock.Any()).Return(proj, nil)
	envRepo.EXPECT().Save(filepath.Join(proj.Path, ".env"), gomock.Any()).Return(nil)
	gitRepo.EXPECT().LoadGlobal().Return(nil)
	gitRepo.EXPECT().UpdateIncludeIf(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("boom"))

	svc := services.NewProjectService(projRepo, envRepo, gitRepo, fs, nil)

	_, err := svc.Create("svc", "/tmp/svc", "", "sh", ".env", domain.EnvVars{}, nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	expectedEnv := filepath.Join(proj.Path, ".env")
	expectedProject := filepath.Join(proj.Path, ".project.hcl")
	if len(fs.removed) != 2 {
		t.Fatalf("expected two removals, got %#v", fs.removed)
	}
	if fs.removed[0] != expectedEnv || fs.removed[1] != expectedProject {
		t.Fatalf("unexpected removal order: %#v", fs.removed)
	}
}
