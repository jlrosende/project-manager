package domain

type Project struct {
	Name      string
	Path      string
	EnvVars   EnvVars
	GitConfig *GitConfig
}
