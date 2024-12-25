package domain

type Project struct {
	Name    string
	Path    string
	EnvVars map[string]string
	Theme   string
}
