package domain

type Project struct {
	Name         string `hcl:"name"`
	Descrption   string `hcl:"description"`
	Path         string
	Shell        string         `hcl:"shell,optional"`
	EnvVarsFile  string         `hcl:"env_vars_file"`
	Environments []*Environment `hcl:"environment,block"`
	EnvVars      EnvVars
}

type Environment struct {
	Name        string `hcl:"name,label"`
	EnvVarsMode string `hcl:"env_vars_mode"`
	EnvVarsFile string `hcl:"env_vars_file"`
	EnvVars     EnvVars
}
