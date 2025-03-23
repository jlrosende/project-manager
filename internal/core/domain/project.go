package domain

import "log/slog"

type Project struct {
	Name         string `hcl:"name"`
	Description  string `hcl:"description"`
	Path         string
	Shell        string         `hcl:"shell,optional"`
	EnvVarsFile  string         `hcl:"env_vars_file"`
	Environments []*Environment `hcl:"environment,block"`
	EnvVars      EnvVars
}

type Environment struct {
	Name        string `hcl:"name,label"`
	Color       string `hcl:"color,optional"`
	EnvVarsMode string `hcl:"env_vars_mode"`
	EnvVarsFile string `hcl:"env_vars_file"`
	EnvVars     EnvVars
}

func (p Project) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", p.Name),
		slog.String("path", p.Path),
		slog.String("shell", p.Shell),
		slog.String("env_vars", p.EnvVarsFile),
	)
}
