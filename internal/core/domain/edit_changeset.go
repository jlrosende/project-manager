package domain

import (
	"errors"
	"fmt"
)

// ChangeSource identifies the origin of a field mutation. It is primarily
// used to surface friendlier error messages when a conflicting input arises.
type ChangeSource string

const (
	// ChangeSourceUnknown indicates the mutation origin was not tracked.
	ChangeSourceUnknown ChangeSource = ""

	// ChangeSourceFlag marks a mutation that originated from an explicit CLI flag.
	ChangeSourceFlag ChangeSource = "flag"

	// ChangeSourceInput marks a mutation loaded from a CLI input file.
	ChangeSourceInput ChangeSource = "cli_input"

	// ChangeSourceSkeleton marks a mutation that originated from a generated
	// skeleton the user modified before persistence.
	ChangeSourceSkeleton ChangeSource = "skeleton"
)

// FieldMutation represents a single field update requested by the user. A nil
// *FieldMutation indicates the field was untouched. An empty string value
// represents an explicit request to clear the field.
type FieldMutation struct {
	Value  any
	Source ChangeSource
}

// ErrMutationType is returned when a mutation cannot be converted to the
// desired concrete type.
var ErrMutationType = errors.New("field mutation has unexpected type")

// StringValue converts the mutation value into a string. Nil values resolve to
// the empty string, which allows callers to treat clears uniformly.
func (m *FieldMutation) StringValue() (string, error) {
	if m == nil {
		return "", nil
	}

	switch v := m.Value.(type) {
	case nil:
		return "", nil
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	case []byte:
		return string(v), nil
	default:
		return "", fmt.Errorf("%w: %T", ErrMutationType, m.Value)
	}
}

// ProjectChangeSet groups all requested project-level mutations. Each map entry
// uses a stable identifier so repositories can translate mutations into
// concrete persistence updates.
type ProjectChangeSet struct {
	Fields map[string]*FieldMutation
}

// EnvironmentChangeSet captures mutations that apply to a specific environment
// within the target project. The Name field identifies the environment.
type EnvironmentChangeSet struct {
	Name   string
	Fields map[string]*FieldMutation
}

// EditChangeSet aggregates project and environment mutations for a single edit
// request. Either scope may be nil when untouched.
type EditChangeSet struct {
	Project     *ProjectChangeSet
	Environment *EnvironmentChangeSet
}

// HasProjectChanges reports whether the change set contains project-level
// mutations.
func (c EditChangeSet) HasProjectChanges() bool {
	return c.Project != nil && len(c.Project.Fields) > 0
}

// HasEnvironmentChanges reports whether the change set contains environment
// mutations.
func (c EditChangeSet) HasEnvironmentChanges() bool {
	return c.Environment != nil && len(c.Environment.Fields) > 0
}

const (
	// ProjectFieldDescription identifies the project description field.
	ProjectFieldDescription = "project.description"
	// ProjectFieldShell identifies the project shell field.
	ProjectFieldShell = "project.shell"
	// ProjectFieldEnvVarsFile identifies the project env vars file field.
	ProjectFieldEnvVarsFile = "project.env_vars_file"
	// ProjectFieldDefaultEnv identifies the project default environment field.
	ProjectFieldDefaultEnv = "project.default_env"
)

const (
	// EnvironmentFieldColor identifies the environment color field.
	EnvironmentFieldColor = "environment.color"
	// EnvironmentFieldEnvVarsMode identifies the env vars mode field.
	EnvironmentFieldEnvVarsMode = "environment.env_vars_mode"
	// EnvironmentFieldEnvVarsFile identifies the environment env vars file field.
	EnvironmentFieldEnvVarsFile = "environment.env_vars_file"
)
