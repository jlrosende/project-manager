package ports

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
