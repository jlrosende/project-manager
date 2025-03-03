package shells

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/creack/pty/v2"
	"golang.org/x/term"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type PseudoShellepository struct {
	cmd  *exec.Cmd
	ptmx *os.File
}

var _ ports.ShellRepository = (*ShellRepository)(nil)

func NewPseudoShellRepository(project *domain.Project, env string, path string) (*PseudoShellepository, error) {

	shellPath, err := exec.LookPath(project.Shell)

	if err != nil {
		return nil, err
	}

	shell := &PseudoShellepository{
		cmd: exec.Command(shellPath),
	}

	// load env vars
	shell.cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("PM_ACTIVE_PROJECT=%s", project.Name),
	)

	if env == "" {
		shell.cmd.Env = append(
			shell.cmd.Env,
			project.EnvVars.ToSlice()...,
		)
	} else {
		for _, e := range project.Environments {
			if e.Name == env {
				if e.EnvVarsMode == domain.ENV_VARS_MODE_MERGE {
					shell.cmd.Env = append(
						shell.cmd.Env,
						project.EnvVars.ToSlice()...,
					)
					shell.cmd.Env = append(
						shell.cmd.Env,
						e.EnvVars.ToSlice()...,
					)
				} else {
					shell.cmd.Env = append(
						shell.cmd.Env,
						e.EnvVars.ToSlice()...,
					)
				}
				break
			}
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if info, err := os.Stat(absPath); err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, errors.New("the path must be a directory")
	}

	if path == "" {
		shell.cmd.Dir = project.Path
	} else {
		shell.cmd.Dir = path
	}

	// set in/out/err
	// shell.cmd.Stdin = os.Stdin
	// shell.cmd.Stdout = os.Stdout
	// shell.cmd.Stderr = os.Stderr

	return shell, nil
}

func (s *PseudoShellepository) Start() (*os.Process, error) {
	ptmx, err := pty.Start(s.cmd)
	if err != nil {
		return nil, err
	}

	s.ptmx = ptmx

	return s.cmd.Process, nil
}

func (s *PseudoShellepository) Wait() (int, error) {

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, s.ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() { _, _ = io.Copy(s.ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, s.ptmx)

	// Test
	// if err := s.cmd.Wait(); err != nil {
	// 	if exiterr, ok := err.(*exec.ExitError); ok {
	// 		return exiterr.ExitCode(), nil
	// 	} else {
	// 		return 0, err
	// 	}
	// }

	return 0, nil
}

func (s *PseudoShellepository) Kill() error {
	return s.cmd.Process.Kill()
}
