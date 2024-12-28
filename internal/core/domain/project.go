package domain

import "fmt"

type Project struct {
	Name    string
	Path    string
	EnvVars map[string]string
	Theme   string
}

func (p Project) EnvVarsSlice() []string {
	envVars := []string{}
	for k, v := range p.EnvVars {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}
	return envVars
}
