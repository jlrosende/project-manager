package bootstrap

import (
	"log/slog"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/ports"
	"github.com/jlrosende/project-manager/internal/core/services"
)

type Options struct {
	Filesystem ports.Filesystem
	Logger     ports.Logger
}

type Container struct {
	Filesystem     ports.Filesystem
	Logger         ports.Logger
	ProjectService ports.ProjectService
	EnvVarsService ports.EnvVarsService
	GitService     ports.GitService
	ProjectRepo    ports.ProjectRepository
	EnvVarsRepo    ports.EnvVarsRepository
	GitRepo        ports.GitRepository
}

func New(opts Options) (*Container, error) {
	fsys := opts.Filesystem
	if fsys == nil {
		fsys = repositories.NewFilesystem()
	}

	logger := opts.Logger

	repoProject, err := repositories.NewProjectRepository(fsys)
	if err != nil {
		return nil, err
	}

	repoEnv, err := repositories.NewEnvVarsRepository()
	if err != nil {
		return nil, err
	}

	repoGit, err := repositories.NewGitRepository()
	if err != nil {
		return nil, err
	}

	projSvc := services.NewProjectService(repoProject, repoEnv, repoGit, fsys, logger)
	envSvc := services.NewEnvVarsServiceService(repoEnv)
	gitSvc := services.NewGitService(repoGit)

	return &Container{
		Filesystem:     fsys,
		Logger:         logger,
		ProjectService: projSvc,
		EnvVarsService: envSvc,
		GitService:     gitSvc,
		ProjectRepo:    repoProject,
		EnvVarsRepo:    repoEnv,
		GitRepo:        repoGit,
	}, nil
}

func WrapSlog(l *slog.Logger) ports.Logger {
	if l == nil {
		return nil
	}

	return repositories.NewLogger(l)
}
