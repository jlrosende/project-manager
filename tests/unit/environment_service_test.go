//go:build unit
// +build unit

package unit_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

func TestEnvironmentService_ApplyAddsNewEnvironment(t *testing.T) {
	t.Helper()

	projectRoot := filepath.Join(t.TempDir(), "demo-project")
	proj := &domain.Project{Name: "demo", Path: projectRoot}
	projectSvc := &fakeProjectService{project: proj}
	envRepo := &fakeEnvVarsRepository{}
	fs := repositories.NewFilesystem()

	svc := services.NewEnvironmentService(projectSvc, envRepo, fs)
	input := &domain.EnvironmentInput{
		Name: ptrString("staging"),
		EnvVars: map[string]string{
			"FOO": "bar",
		},
	}

	if err := svc.Apply("demo", input, services.EnvironmentApplyOptions{}); err != nil {
		t.Fatalf("apply environment: %v", err)
	}

	if len(projectSvc.addCalls) != 1 {
		t.Fatalf("expected one AddEnvironment call, got %d", len(projectSvc.addCalls))
	}

	call := projectSvc.addCalls[0]
	if call.env.Name != "staging" {
		t.Fatalf("expected environment name staging, got %s", call.env.Name)
	}

	if call.env.EnvVarsFile != ".staging.env" {
		t.Fatalf("expected default env vars file .staging.env, got %s", call.env.EnvVarsFile)
	}

	if call.env.EnvVarsMode != domain.EnvVarsModeMerge {
		t.Fatalf("expected env vars mode merge, got %s", call.env.EnvVarsMode)
	}

	if call.env.Color != "" {
		t.Fatalf("expected empty color, got %s", call.env.Color)
	}

	if call.envVars["FOO"] != "bar" {
		t.Fatalf("expected env vars map to contain FOO=bar, got %v", call.envVars)
	}

	if len(envRepo.saves) != 0 {
		t.Fatalf("expected env repo Save not to be called for AddEnvironment")
	}
}

func TestEnvironmentService_ApplyForceUpdatesEnvironment(t *testing.T) {
	t.Helper()

	projectRoot := filepath.Join(t.TempDir(), "demo-project")
	proj := &domain.Project{
		Name: "demo",
		Path: projectRoot,
		Environments: []*domain.Environment{
			{Name: "staging", EnvVarsFile: ".staging.env", EnvVarsMode: domain.EnvVarsModeMerge, Color: "42"},
		},
	}

	projectSvc := &fakeProjectService{project: proj}
	envRepo := &fakeEnvVarsRepository{}
	fs := repositories.NewFilesystem()

	svc := services.NewEnvironmentService(projectSvc, envRepo, fs)
	input := &domain.EnvironmentInput{
		Name:        ptrString("staging"),
		EnvVarsFile: ptrString(".env.updated"),
		EnvVarsMode: ptrString("replace"),
		Color:       ptrString("200"),
		EnvVars: map[string]string{
			"FORCE": "1",
		},
	}

	if err := svc.Apply("demo", input, services.EnvironmentApplyOptions{Force: true}); err != nil {
		t.Fatalf("force apply environment: %v", err)
	}

	if len(projectSvc.updateCalls) != 1 {
		t.Fatalf("expected one UpdateEnvironment call, got %d", len(projectSvc.updateCalls))
	}

	call := projectSvc.updateCalls[0]
	if call.original != "staging" {
		t.Fatalf("expected original env name staging, got %s", call.original)
	}

	updated := call.env
	if updated.EnvVarsFile != ".env.updated" {
		t.Fatalf("expected env vars file .env.updated, got %s", updated.EnvVarsFile)
	}

	if updated.EnvVarsMode != domain.EnvVarsModeReplace {
		t.Fatalf("expected env vars mode replace, got %s", updated.EnvVarsMode)
	}

	if updated.Color != "200" {
		t.Fatalf("expected color 200, got %s", updated.Color)
	}

	if len(envRepo.saves) != 1 {
		t.Fatalf("expected env repo Save to be called once, got %d", len(envRepo.saves))
	}

	save := envRepo.saves[0]
	expectedPath := filepath.Join(projectRoot, ".env.updated")
	if save.path != expectedPath {
		t.Fatalf("expected env vars to be saved at %s, got %s", expectedPath, save.path)
	}

	if save.vars["FORCE"] != "1" {
		t.Fatalf("expected env vars to contain FORCE=1, got %v", save.vars)
	}
}

func TestEnvironmentService_ApplyExistingWithoutForceFails(t *testing.T) {
	t.Helper()

	projectRoot := filepath.Join(t.TempDir(), "demo-project")
	proj := &domain.Project{
		Name: "demo",
		Path: projectRoot,
		Environments: []*domain.Environment{
			{Name: "staging", EnvVarsFile: ".staging.env", EnvVarsMode: domain.EnvVarsModeMerge},
		},
	}

	projectSvc := &fakeProjectService{project: proj}
	envRepo := &fakeEnvVarsRepository{}
	fs := repositories.NewFilesystem()

	svc := services.NewEnvironmentService(projectSvc, envRepo, fs)
	input := &domain.EnvironmentInput{
		Name: ptrString("staging"),
	}

	err := svc.Apply("demo", input, services.EnvironmentApplyOptions{})
	if err == nil {
		t.Fatalf("expected error when environment exists without force")
	}

	if !errors.Is(err, services.ErrEnvironmentExists) {
		t.Fatalf("expected ErrEnvironmentExists, got %v", err)
	}

	if len(projectSvc.updateCalls) != 0 {
		t.Fatalf("expected no UpdateEnvironment calls")
	}
}

func TestEnvironmentService_InvalidModeFails(t *testing.T) {
	t.Helper()

	projectRoot := filepath.Join(t.TempDir(), "demo-project")
	projectSvc := &fakeProjectService{project: &domain.Project{Name: "demo", Path: projectRoot}}
	envRepo := &fakeEnvVarsRepository{}
	fs := repositories.NewFilesystem()

	svc := services.NewEnvironmentService(projectSvc, envRepo, fs)
	input := &domain.EnvironmentInput{
		Name:        ptrString("staging"),
		EnvVarsMode: ptrString("invalid"),
	}

	if err := svc.Apply("demo", input, services.EnvironmentApplyOptions{}); err == nil {
		t.Fatalf("expected validation error for invalid env vars mode")
	}
}

type fakeProjectService struct {
	project     *domain.Project
	addCalls    []addCall
	updateCalls []updateCall
}

type addCall struct {
	env     domain.Environment
	envVars domain.EnvVars
}

type updateCall struct {
	original string
	env      *domain.Environment
}

func (f *fakeProjectService) Load(name string) (*domain.Project, error) {
	return f.project, nil
}

func (f *fakeProjectService) List() ([]*domain.Project, error) { return nil, nil }

func (f *fakeProjectService) Probe(name string) (domain.ProjectExistence, error) {
	return domain.ProjectExistence{}, nil
}

func (f *fakeProjectService) Create(name, path, subproject, shell, envFile string, envVars domain.EnvVars, git *domain.GitConfig) (*domain.Project, error) {
	return nil, nil
}

func (f *fakeProjectService) AddEnvironment(projectName string, env *domain.Environment, envVars domain.EnvVars) error {
	copyEnv := *env
	copyVars := domain.EnvVars{}
	for k, v := range envVars {
		copyVars[k] = v
	}

	f.addCalls = append(f.addCalls, addCall{env: copyEnv, envVars: copyVars})
	return nil
}

func (f *fakeProjectService) UpdateProject(project *domain.Project) error { return nil }

func (f *fakeProjectService) UpdateEnvironment(projectName, originalEnvName string, env *domain.Environment) error {
	f.updateCalls = append(f.updateCalls, updateCall{original: originalEnvName, env: env})
	return nil
}

func (f *fakeProjectService) Delete(name string) error { return nil }

func (f *fakeProjectService) DeleteProject(ctx context.Context, options domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
	return nil, nil
}

type fakeEnvVarsRepository struct {
	saves []saveCall
}

type saveCall struct {
	path string
	vars map[string]string
}

func (f *fakeEnvVarsRepository) Load(path string) (domain.EnvVars, error) {
	return nil, nil
}

func (f *fakeEnvVarsRepository) Save(path string, envVars map[string]string) error {
	copyVars := make(map[string]string, len(envVars))
	for k, v := range envVars {
		copyVars[k] = v
	}

	f.saves = append(f.saves, saveCall{path: path, vars: copyVars})
	return nil
}

func (f *fakeEnvVarsRepository) Delete(ctx context.Context, path string) error { return nil }

func ptrString(value string) *string {
	v := value
	return &v
}
