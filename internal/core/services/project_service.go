package services

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
	"github.com/jlrosende/project-manager/internal/tools"
)

type ProjectService struct {
	git     ports.GitRepository
	envVars ports.EnvVarsRepository
	project ports.ProjectRepository
}

var _ ports.ProjectService = (*ProjectService)(nil)

func NewProjectService(
	project ports.ProjectRepository,
	envVars ports.EnvVarsRepository,
	git ports.GitRepository,
) *ProjectService {
	return &ProjectService{
		project: project,
		envVars: envVars,
		git:     git,
	}
}

func (svc *ProjectService) Load(name string) (*domain.Project, error) {
	slog.Debug("load project", slog.String("name", name))

	projects, err := svc.project.List()
	if err != nil {
		slog.Error("list projects failed", slog.String("err", err.Error()))
		return nil, err
	}

	for _, project := range projects {
		if project.Name != name {
			continue
		}

		var envVarsPath string
		if filepath.IsAbs(project.EnvVarsFile) {
			envVarsPath = project.EnvVarsFile
		} else {
			envVarsPath = filepath.Join(project.Path, project.EnvVarsFile)
		}

		project.EnvVars, err = svc.envVars.Load(envVarsPath)
		if err != nil {
			slog.Error(
				"load project env vars failed",
				slog.String("path", envVarsPath),
				slog.String("err", err.Error()),
			)

			return nil, err
		}

		for _, env := range project.Environments {
			var envVarsPath string
			if filepath.IsAbs(env.EnvVarsFile) {
				envVarsPath = env.EnvVarsFile
			} else {
				envVarsPath = filepath.Join(project.Path, env.EnvVarsFile)
			}

			env.EnvVars, err = svc.envVars.Load(envVarsPath)
			if err != nil {
				slog.Warn(
					"env vars load failed; using empty",
					slog.String("env", env.Name),
					slog.String("path", envVarsPath),
				)
				env.EnvVars = domain.EnvVars{}
			}
		}

		slog.Debug("loaded project", slog.String("name", project.Name))

		return project, nil
	}

	slog.Warn("project not found", slog.String("name", name))

	return nil, fmt.Errorf("project '%s' not found, projects are case sensitive", name)
}

func (svc *ProjectService) List() ([]*domain.Project, error) {
	return svc.project.List()
}

func (svc *ProjectService) Create(
	name, path, subproject, shell, envFile string,
	envVars domain.EnvVars,
	gitConfig *domain.GitConfig,
) (*domain.Project, error) {
	slog.Info(
		"create project",
		slog.String("name", name),
		slog.String("path", path),
		slog.String("subproject", subproject),
	)

	file := strings.TrimSpace(envFile)
	if file == "" {
		file = ".env"
	}

	proj, err := svc.project.Create(name, path, subproject, shell, file, envVars, gitConfig)
	if err != nil {
		slog.Error("project create failed", slog.String("err", err.Error()))
		return nil, err
	}

	if proj != nil {
		created := []string{}

		p := file
		if !tools.IsAbs(p) {
			p = filepath.Join(proj.Path, p)
		}

		if err := svc.envVars.Save(p, envVars); err != nil {
			slog.Error("env file save failed", slog.String("path", p), slog.String("err", err.Error()))

			_ = os.Remove(filepath.Join(proj.Path, ".project.hcl"))

			return nil, err
		}

		created = append(created, p)

		if err := svc.git.LoadGlobal(); err != nil {
			slog.Error("git load global failed", slog.String("err", err.Error()))

			_ = os.Remove(p)
			_ = os.Remove(filepath.Join(proj.Path, ".project.hcl"))

			return nil, err
		}

		gitdir := fmt.Sprintf("gitdir/i:%s/", proj.Path)

		perProj := filepath.Join(proj.Path, fmt.Sprintf(".%s.gitconfig", proj.Name))
		if err := svc.git.UpdateIncludeIf(gitdir, perProj, subproject); err != nil {
			slog.Error(
				"git includeIf update failed",
				slog.String("gitdir", gitdir),
				slog.String("path", perProj),
				slog.String("err", err.Error()),
			)

			_ = os.Remove(p)
			_ = os.Remove(filepath.Join(proj.Path, ".project.hcl"))

			return nil, err
		}

		if home, e := os.UserHomeDir(); e != nil || svc.git.SaveGlobal(home) != nil {
			slog.Error("git save global failed")

			_ = os.Remove(p)
			_ = os.Remove(filepath.Join(proj.Path, ".project.hcl"))

			return nil, fmt.Errorf("failed to save global git config")
		}

		if gitConfig != nil {
			if err := svc.git.Save(perProj, gitConfig); err != nil {
				slog.Error("per-project git save failed", slog.String("path", perProj), slog.String("err", err.Error()))

				_ = os.Remove(p)
				_ = os.Remove(filepath.Join(proj.Path, ".project.hcl"))

				return nil, err
			}

			created = append(created, perProj)
		}

		_ = created

		slog.Info("project created", slog.String("name", proj.Name))
	}

	return proj, nil
}

func (svc *ProjectService) AddEnvironment(projectName string, env *domain.Environment, envVars domain.EnvVars) error {
	slog.Info("add environment", slog.String("project", projectName), slog.String("env", func() string {
		if env != nil {
			return env.Name
		}

		return ""
	}()))

	if env == nil {
		slog.Error("add environment failed: nil env")
		return fmt.Errorf("env is nil")
	}

	if strings.TrimSpace(env.EnvVarsMode) == "" {
		env.EnvVarsMode = domain.EnvVarsModeMerge
	}

	if strings.TrimSpace(env.EnvVarsFile) == "" {
		slug := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(env.Name, " ", "-")))
		if slug != "" {
			env.EnvVarsFile = "." + slug + ".env"
		} else {
			env.EnvVarsFile = ".env"
		}
	}

	proj, err := svc.Load(projectName)
	if err != nil {
		slog.Error("load project failed", slog.String("project", projectName), slog.String("err", err.Error()))
		return err
	}

	for _, e := range proj.Environments {
		if e.Name == env.Name {
			slog.Warn("environment exists", slog.String("project", projectName), slog.String("env", env.Name))
			return fmt.Errorf("environment %s already exists", env.Name)
		}
	}

	p := env.EnvVarsFile
	if !tools.IsAbs(p) {
		p = filepath.Join(proj.Path, p)
	}

	if err := svc.envVars.Save(p, envVars); err != nil {
		slog.Error("env file save failed", slog.String("path", p), slog.String("err", err.Error()))
		return err
	}

	if err := svc.project.AddEnvironment(projectName, env, envVars); err != nil {
		slog.Error(
			"add environment repo failed",
			slog.String("project", projectName),
			slog.String("env", env.Name),
			slog.String("err", err.Error()),
		)

		_ = os.Remove(p)

		return err
	}

	slog.Info("environment added", slog.String("project", projectName), slog.String("env", env.Name))

	return nil
}

func (svc *ProjectService) UpdateProject(project *domain.Project) error {
	slog.Info("update project", slog.String("name", project.Name))

	err := svc.project.UpdateProject(project)
	if err != nil {
		slog.Error("update project failed", slog.String("name", project.Name), slog.String("err", err.Error()))
		return err
	}

	slog.Info("project updated", slog.String("name", project.Name))

	return nil
}

func (svc *ProjectService) UpdateEnvironment(projectName, originalEnvName string, env *domain.Environment) error {
	slog.Info(
		"update environment",
		slog.String("project", projectName),
		slog.String("original", originalEnvName),
		slog.String("env", func() string {
			if env != nil {
				return env.Name
			}

			return ""
		}()),
	)

	err := svc.project.UpdateEnvironment(projectName, originalEnvName, env)
	if err != nil {
		slog.Error(
			"update environment failed",
			slog.String("project", projectName),
			slog.String("original", originalEnvName),
			slog.String("err", err.Error()),
		)

		return err
	}

	slog.Info("environment updated", slog.String("project", projectName), slog.String("env", func() string {
		if env != nil {
			return env.Name
		}

		return ""
	}()))

	return nil
}

func (svc *ProjectService) Delete(name string) error {
	slog.Info("delete project", slog.String("name", name))

	err := svc.project.Delete(name)
	if err != nil {
		slog.Error("delete project failed", slog.String("name", name), slog.String("err", err.Error()))
		return err
	}

	slog.Info("project deleted", slog.String("name", name))

	return nil
}
