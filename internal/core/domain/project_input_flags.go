package domain

// ProjectConfigFlags captures CLI-provided overrides for project creation.
type ProjectConfigFlags struct {
	Name           string
	Path           string
	Here           bool
	Description    string
	Shell          string
	EnvFile        string
	PathSet        bool
	HereSet        bool
	DescriptionSet bool
	ShellSet       bool
	EnvFileSet     bool
	Environment    EnvironmentConfigFlags
}

// EnvironmentConfigFlags captures CLI-provided overrides for environment creation.
type EnvironmentConfigFlags struct {
	Name        string
	NameSet     bool
	EnvVarsFile string
	EnvFileSet  bool
	EnvVarsMode string
	ModeSet     bool
	Color       string
	ColorSet    bool
	EnvVars     map[string]string
}

// HasInput reports whether any environment-related flag values were provided.
func (f EnvironmentConfigFlags) HasInput() bool {
	return f.NameSet || f.EnvFileSet || f.ModeSet || f.ColorSet || len(f.EnvVars) > 0
}

// ProjectCreateOptions configures filesystem creation behaviour.
type ProjectCreateOptions struct {
	Force bool
}

// EnvironmentApplyOptions configures behaviour when applying environment inputs.
type EnvironmentApplyOptions struct {
	Force bool
}
