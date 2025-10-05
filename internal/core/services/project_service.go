package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type ProjectService struct {
	git     ports.GitRepository
	envVars ports.EnvVarsRepository
	project ports.ProjectRepository
	fs      ports.Filesystem
	logger  ports.Logger
}

var _ ports.ProjectService = (*ProjectService)(nil)

func NewProjectService(
	project ports.ProjectRepository,
	envVars ports.EnvVarsRepository,
	git ports.GitRepository,
	fs ports.Filesystem,
	logger ports.Logger,
) *ProjectService {
	if project == nil || envVars == nil || git == nil || fs == nil {
		panic("project service dependencies must not be nil")
	}

	return &ProjectService{
		project: project,
		envVars: envVars,
		git:     git,
		fs:      fs,
		logger:  logger,
	}
}

func (svc *ProjectService) Load(name string) (*domain.Project, error) {
	svc.logDebug("load project", field("name", name))

	projects, err := svc.project.List()
	if err != nil {
		svc.logError("list projects failed", field("err", err))
		return nil, err
	}

	for _, project := range projects {
		if project.Name != name {
			continue
		}

		envVarsPath := project.EnvVarsFile
		if !svc.fs.IsAbs(envVarsPath) {
			envVarsPath = svc.fs.Join(project.Path, envVarsPath)
		}

		project.EnvVars, err = svc.envVars.Load(envVarsPath)
		if err != nil {
			svc.logError("load project env vars failed", field("path", envVarsPath), field("err", err))
			return nil, err
		}

		for _, env := range project.Environments {
			envVarsPath := env.EnvVarsFile
			if !svc.fs.IsAbs(envVarsPath) {
				envVarsPath = svc.fs.Join(project.Path, envVarsPath)
			}

			env.EnvVars, err = svc.envVars.Load(envVarsPath)
			if err != nil {
				svc.logWarn(
					"env vars load failed; using empty",
					field("env", env.Name),
					field("path", envVarsPath),
				)
				env.EnvVars = domain.EnvVars{}
			}
		}

		svc.logDebug("loaded project", field("name", project.Name))

		return project, nil
	}

	svc.logWarn("project not found", field("name", name))

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
	svc.logInfo(
		"create project",
		field("name", name),
		field("path", path),
		field("subproject", subproject),
	)

	file := strings.TrimSpace(envFile)
	if file == "" {
		file = ".env"
	}

	proj, err := svc.project.Create(name, path, subproject, shell, file, envVars, gitConfig)
	if err != nil {
		svc.logError("project create failed", field("err", err))
		return nil, err
	}

	if proj != nil {
		projectFile := svc.fs.Join(proj.Path, ".project.hcl")

		p := file
		if !svc.fs.IsAbs(p) {
			p = svc.fs.Join(proj.Path, p)
		}

		if err := svc.envVars.Save(p, envVars); err != nil {
			svc.logError("env file save failed", field("path", p), field("err", err))
			_ = svc.fs.Remove(projectFile)

			return nil, err
		}

		if err := svc.git.LoadGlobal(); err != nil {
			svc.logError("git load global failed", field("err", err))
			_ = svc.fs.Remove(p)
			_ = svc.fs.Remove(projectFile)

			return nil, err
		}

		gitdir := fmt.Sprintf("gitdir/i:%s/", proj.Path)
		perProj := svc.fs.Join(proj.Path, fmt.Sprintf(".%s.gitconfig", proj.Name))

		if err := svc.git.UpdateIncludeIf(gitdir, perProj, subproject); err != nil {
			svc.logError(
				"git includeIf update failed",
				field("gitdir", gitdir),
				field("path", perProj),
				field("err", err),
			)
			_ = svc.fs.Remove(p)
			_ = svc.fs.Remove(projectFile)

			return nil, err
		}

		home, homeErr := svc.fs.UserHomeDir()
		if homeErr != nil {
			svc.logError("resolve user home failed", field("err", homeErr))
			_ = svc.fs.Remove(p)
			_ = svc.fs.Remove(projectFile)

			return nil, fmt.Errorf("failed to resolve user home directory")
		}

		if err := svc.git.SaveGlobal(home); err != nil {
			svc.logError("git save global failed", field("err", err))
			_ = svc.fs.Remove(p)
			_ = svc.fs.Remove(projectFile)

			return nil, fmt.Errorf("failed to save global git config")
		}

		if gitConfig != nil {
			if err := svc.git.Save(perProj, gitConfig); err != nil {
				svc.logError("per-project git save failed", field("path", perProj), field("err", err))
				_ = svc.fs.Remove(p)
				_ = svc.fs.Remove(projectFile)

				return nil, err
			}
		} else {
			if err := svc.fs.WriteFile(perProj, []byte{}, 0o600); err != nil {
				svc.logError("per-project git file create failed", field("path", perProj), field("err", err))
				_ = svc.fs.Remove(p)
				_ = svc.fs.Remove(projectFile)

				return nil, err
			}
		}

		svc.logInfo("project created", field("name", proj.Name))
	}

	return proj, nil
}

func (svc *ProjectService) AddEnvironment(projectName string, env *domain.Environment, envVars domain.EnvVars) error {
	svc.logInfo(
		"add environment",
		field("project", projectName),
		field("env", func() string {
			if env != nil {
				return env.Name
			}

			return ""
		}()),
	)

	if env == nil {
		svc.logError("add environment failed: nil env")
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
		svc.logError("load project failed", field("project", projectName), field("err", err))
		return err
	}

	for _, e := range proj.Environments {
		if e.Name == env.Name {
			svc.logWarn("environment exists", field("project", projectName), field("env", env.Name))
			return fmt.Errorf("environment %s already exists", env.Name)
		}
	}

	p := env.EnvVarsFile
	if !svc.fs.IsAbs(p) {
		p = svc.fs.Join(proj.Path, p)
	}

	if err := svc.envVars.Save(p, envVars); err != nil {
		svc.logError("env file save failed", field("path", p), field("err", err))
		return err
	}

	if err := svc.project.AddEnvironment(projectName, env, envVars); err != nil {
		svc.logError(
			"add environment repo failed",
			field("project", projectName),
			field("env", env.Name),
			field("err", err),
		)
		_ = svc.fs.Remove(p)

		return err
	}

	svc.logInfo("environment added", field("project", projectName), field("env", env.Name))

	return nil
}

func (svc *ProjectService) UpdateProject(project *domain.Project) error {
	svc.logInfo("update project", field("name", project.Name))

	err := svc.project.UpdateProject(project)
	if err != nil {
		svc.logError("update project failed", field("name", project.Name), field("err", err))
		return err
	}

	svc.logInfo("project updated", field("name", project.Name))

	return nil
}

func (svc *ProjectService) UpdateEnvironment(projectName, originalEnvName string, env *domain.Environment) error {
	svc.logInfo(
		"update environment",
		field("project", projectName),
		field("original", originalEnvName),
		field("env", func() string {
			if env != nil {
				return env.Name
			}

			return ""
		}()),
	)

	err := svc.project.UpdateEnvironment(projectName, originalEnvName, env)
	if err != nil {
		svc.logError(
			"update environment failed",
			field("project", projectName),
			field("original", originalEnvName),
			field("err", err),
		)

		return err
	}

	svc.logInfo("environment updated", field("project", projectName), field("env", func() string {
		if env != nil {
			return env.Name
		}

		return ""
	}()))

	return nil
}

func (svc *ProjectService) Delete(name string) error {
	svc.logInfo("delete project", field("name", name))

	err := svc.project.Delete(name)
	if err != nil {
		svc.logError("delete project failed", field("name", name), field("err", err))
		return err
	}

	svc.logInfo("project deleted", field("name", name))

	return nil
}

func (svc *ProjectService) DeleteProject(ctx context.Context, options domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}

	identifier, err := svc.project.ResolveIdentifier(ctx, options.Target)
	if err != nil {
		svc.logError("resolve delete target failed", field("target", options.Target.Name), field("err", err))
		return nil, err
	}

	if identifier.Status.Locked {
		svc.logWarn("project locked", field("name", identifier.Name))
		return nil, domain.ErrProjectLocked
	}

	if identifier.Status.ReadOnly {
		svc.logWarn("workspace read-only", field("name", identifier.Name))
		return nil, domain.ErrWorkspaceReadOnly
	}

	plan, err := svc.fs.PlanDeletion(ctx, identifier, options.Scope, options.Backup)
	if err != nil {
		svc.logError("plan project deletion failed", field("name", identifier.Name), field("err", err))
		return nil, fmt.Errorf("plan deletion: %w", err)
	}

	result := &domain.ProjectDeleteResult{
		Scope: options.Scope,
		Plan:  plan,
	}

	if options.DryRun && options.RequiresBackup() && options.Backup != nil {
		result.BackupPath = options.Backup.Destination
	}

	if options.DryRun {
		svc.logInfo(
			"delete preview",
			field("name", identifier.Name),
			field("scope", plan.Scope.String()),
			field("artifacts", len(plan.Artifacts)),
		)
		return result, nil
	}

	var (
		removed []domain.DeletionArtifact
		skipped []domain.DeletionArtifact
		errs    []error
	)

	hasGitArtifacts := false
	for _, artifact := range plan.Artifacts {
		if artifact.Type == domain.ArtifactHook || artifact.Type == domain.ArtifactGitInclude {
			hasGitArtifacts = true
			break
		}
	}

	var backupMeta *domain.BackupArtifact
	if options.RequiresBackup() {
		backupMeta, err = svc.fs.CreateBackup(ctx, identifier, options.Backup)
		if err != nil {
			svc.logError("create backup failed", field("name", identifier.Name), field("err", err))
			return nil, fmt.Errorf("create backup: %w", err)
		}

		if backupMeta != nil && backupMeta.Completed() {
			result.BackupPath = backupMeta.FinalPath
			svc.logInfo("backup created", field("path", backupMeta.FinalPath))
		}
	}

	failureIndex := -1
	gitCleaned := false

	if !hasGitArtifacts && plan.Scope.IncludesMetadata() {
		if err := svc.git.RemoveHooks(ctx, identifier); err != nil {
			errs = append(errs, fmt.Errorf("remove git hooks: %w", err))
			skipped = append(skipped, domain.DeletionArtifact{
				Type:        domain.ArtifactGitInclude,
				Path:        identifier.Path,
				Description: "git hooks",
			})

			result.ArtifactsRemoved = removed
			result.ArtifactsSkipped = skipped
			result.Errors = errs

			return result, fmt.Errorf("delete project: %w", domain.ErrPartialDeletion)
		}

		svc.logInfo("git hooks removed", field("name", identifier.Name))
		gitCleaned = true
	}

	for idx, artifact := range plan.Artifacts {
		svc.logInfo("delete artifact", field("type", artifact.Type), field("path", artifact.Path))

		switch artifact.Type {
		case domain.ArtifactEnvVars:
			if err := svc.envVars.Delete(ctx, artifact.Path); err != nil {
				errs = append(errs, fmt.Errorf("remove env vars: %w", err))
				skipped = append(skipped, artifact)
				failureIndex = idx
			}
			if failureIndex == -1 {
				removed = append(removed, artifact)
			}
		case domain.ArtifactHook, domain.ArtifactGitInclude:
			if gitCleaned {
				removed = append(removed, artifact)
				continue
			}

			if err := svc.git.RemoveHooks(ctx, identifier); err != nil {
				errs = append(errs, fmt.Errorf("remove git hooks: %w", err))
				skipped = append(skipped, artifact)
				failureIndex = idx
			} else {
				removed = append(removed, artifact)
				gitCleaned = true
			}
		case domain.ArtifactBackup:
			if options.RequiresBackup() {
				if backupMeta != nil && backupMeta.Completed() {
					removed = append(removed, artifact)
				} else {
					errs = append(errs, fmt.Errorf("backup artifact not created"))
					skipped = append(skipped, artifact)
					failureIndex = idx
				}
			}
		default:
			// handled by filesystem execution below
		}

		if failureIndex != -1 {
			for _, pending := range plan.Artifacts[idx+1:] {
				skipped = append(skipped, pending)
			}
			break
		}
	}

	if failureIndex == -1 {
		fsRemoved, execErr := svc.fs.ExecuteDeletion(ctx, plan)
		if execErr != nil {
			errs = append(errs, fmt.Errorf("execute deletion: %w", execErr))
		} else if len(fsRemoved) > 0 {
			for _, artifact := range fsRemoved {
				if artifact.Type == domain.ArtifactRegistry {
					continue
				}
				removed = append(removed, artifact)
			}
		}

		if len(errs) == 0 {
			if err := svc.project.FinalizeDeletion(ctx, identifier, plan.Scope); err != nil {
				errs = append(errs, fmt.Errorf("finalize deletion: %w", err))
			} else {
				svc.logInfo("project deletion finalized", field("name", identifier.Name))
			}
		}
	}

	result.ArtifactsRemoved = removed
	result.ArtifactsSkipped = skipped
	result.Errors = errs

	if len(errs) > 0 {
		svc.logWarn("project deletion completed with errors", field("name", identifier.Name), field("errors", len(errs)))
		return result, fmt.Errorf("delete project: %w", domain.ErrPartialDeletion)
	}

	svc.logInfo("project deletion complete", field("name", identifier.Name), field("scope", plan.Scope.String()))

	return result, nil
}

func (svc *ProjectService) logDebug(msg string, fields ...ports.LogField) {
	if svc.logger != nil {
		svc.logger.Debug(msg, fields...)
	}
}

func (svc *ProjectService) logInfo(msg string, fields ...ports.LogField) {
	if svc.logger != nil {
		svc.logger.Info(msg, fields...)
	}
}

func (svc *ProjectService) logWarn(msg string, fields ...ports.LogField) {
	if svc.logger != nil {
		svc.logger.Warn(msg, fields...)
	}
}

func (svc *ProjectService) logError(msg string, fields ...ports.LogField) {
	if svc.logger != nil {
		svc.logger.Error(msg, fields...)
	}
}

func field(key string, value any) ports.LogField {
	return ports.LogField{Key: key, Value: value}
}
