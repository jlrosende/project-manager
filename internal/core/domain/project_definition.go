package domain

import "strings"

// ProjectDefinition captures the normalized inputs required to create a project.
type ProjectDefinition struct {
	Name        string
	Here        bool
	Path        string
	Environment *EnvironmentInput
	Metadata    map[string]string
}

// DestinationProvided reports whether either Here or Path has been specified.
func (p ProjectDefinition) DestinationProvided() bool {
	if p.Here {
		return true
	}

	return strings.TrimSpace(p.Path) != ""
}

// ConfigInput models the optional configuration supplied via JSON or YAML files.
// Pointer fields differentiate between omitted and zero-valued inputs for later merge logic.
type ConfigInput struct {
	Name               *string           `json:"name"        yaml:"name"`
	Here               *bool             `json:"here"        yaml:"here"`
	Path               *string           `json:"path"        yaml:"path"`
	Environment        *EnvironmentInput `json:"environment" yaml:"environment"`
	Metadata           map[string]string `json:"metadata"    yaml:"metadata"`
	unknown            map[string]any
	legacyEnvironments bool
}

// EnvironmentInput captures environment-specific configuration supplied either via
// CLI input files or merged flag data.
type EnvironmentInput struct {
	Name        *string           `json:"name"          yaml:"name"`
	EnvVarsFile *string           `json:"env_vars_file" yaml:"env_vars_file"`
	EnvVarsMode *string           `json:"env_vars_mode" yaml:"env_vars_mode"`
	Color       *string           `json:"color"         yaml:"color"`
	EnvVars     map[string]string `json:"env_vars"      yaml:"env_vars"`
}

// UnknownFields returns any extra fields encountered during parsing.
// It returns nil when no unknown fields were recorded.
func (c ConfigInput) UnknownFields() map[string]any {
	return c.unknown
}

// SetUnknownFields records unknown keys captured by config parsing logic.
func (c *ConfigInput) SetUnknownFields(fields map[string]any) {
	if len(fields) == 0 {
		c.unknown = nil
		return
	}

	c.unknown = make(map[string]any, len(fields))
	for k, v := range fields {
		c.unknown[k] = v
	}
}

// HasLegacyEnvironments reports whether the parsed config contained the deprecated
// `environments` map.
func (c ConfigInput) HasLegacyEnvironments() bool {
	return c.legacyEnvironments
}

// SetLegacyEnvironments flags whether the deprecated `environments` map was observed during parsing.
func (c *ConfigInput) SetLegacyEnvironments(flag bool) {
	c.legacyEnvironments = flag
}
