package domain

type Project struct {
	Name         string        `hcl:"name"`
	Descrption   string        `hcl:"description"`
	EnvVarsFile  string        `hcl:"env_vars_file"`
	Environments []Environment `hcl:"environment,block"`
}

type Environment struct {
	Name        string `hcl:"name,label"`
	EnvVarsMode string `hcl:"env_vars_mode"`
	EnvVarsFile string `hcl:"env_vars_file"`
}
