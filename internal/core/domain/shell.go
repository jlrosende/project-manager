package domain

type Shell struct {
	Dir     string
	EnvVars map[string]string
}

func NewShell(shell string) *Shell {
	return &Shell{}
}
