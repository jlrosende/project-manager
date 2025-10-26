package ports

//go:generate go tool mockgen -source=shell_port.go -destination=../../../mocks/mock_shell_port.go -package=mocks

import "os"

type ShellService interface {
	Start() (*os.Process, error)
	Wait() (int, error)
	Kill() error
}

type ShellRepository interface {
	Start() (*os.Process, error)
	Wait() (int, error)
	Kill() error
}
