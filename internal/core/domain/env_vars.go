package domain

import "fmt"

type EnvVars map[string]string

func (e EnvVars) ToSlice() []string {
	envVars := []string{}
	for k, v := range e {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}
	return envVars
}
