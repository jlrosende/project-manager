package domain

type ShellType string

const (
	SHELL_BASH   ShellType = "bash"
	SHELL_SH     ShellType = "sh"
	SHELL_ZSH    ShellType = "zsh"
	SHELL_FISH   ShellType = "fish"
	SHELL_ELVISH ShellType = "elvish"
	SHELL_PWSH   ShellType = "pwsh"
)

type Shell struct {
	Type    ShellType
	Dir     string
	EnvVars map[string]string
}

func NewShell(shell string) *Shell {
	var shellType ShellType
	switch shell {
	default:
		shellType = SHELL_BASH
	}
	return &Shell{
		Type: shellType,
	}
}
