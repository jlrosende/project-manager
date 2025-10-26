package ports

import "github.com/jlrosende/project-manager/internal/core/domain"

// ProjectCreator exposes project creation workflows to adapters.
type ProjectCreator interface {
	Create(
		def domain.ProjectDefinition,
		envVars domain.EnvVars,
		opts domain.ProjectCreateOptions,
	) (*domain.Project, error)
}

// EnvironmentManager exposes environment add/update workflows to adapters.
type EnvironmentManager interface {
	Apply(projectName string, input *domain.EnvironmentInput, opts domain.EnvironmentApplyOptions) error
}

// ProjectInputService orchestrates configuration parsing and merging for adapters.
type ProjectInputService interface {
	LoadConfig(path string) (*domain.ConfigInput, error)
	MergeInputs(
		cfg *domain.ConfigInput,
		flags domain.ProjectConfigFlags,
	) (domain.ProjectDefinition, domain.EnvVars, error)
	ParseEnvVarFlags(pairs []string, flagLabel string) (map[string]string, error)
	EnvironmentProvided(input *domain.EnvironmentInput) bool
}
