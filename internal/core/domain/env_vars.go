package domain

import "fmt"

const (
	ENV_VARS_MODE_MERGE   string = "merge"
	ENV_VARS_MODE_REPLACE string = "replace"
)

type EnvVars map[string]string

func (e EnvVars) ToSlice() []string {
	envVars := []string{}
	for k, v := range e {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}
	return envVars
}
