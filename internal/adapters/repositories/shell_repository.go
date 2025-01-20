package repositories

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type ShellRepository struct {
	cmd *exec.Cmd
}

var _ ports.ShellRepository = (*ShellRepository)(nil)

func NewShellRepository(project *domain.Project, env string, path string) (*ShellRepository, error) {

	shellPath, err := exec.LookPath(project.Shell)

	if err != nil {
		return nil, err
	}

	shell := &ShellRepository{
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
	shell.cmd.Stdin = os.Stdin
	shell.cmd.Stdout = os.Stdout
	shell.cmd.Stderr = os.Stderr

	return shell, nil
}

func (s *ShellRepository) Start() (*os.Process, error) {

	if err := s.cmd.Start(); err != nil {
		return nil, err
	}

	return s.cmd.Process, nil
}

func (s *ShellRepository) Wait() (int, error) {
	if err := s.cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode(), nil
		} else {
			return 0, err
		}
	}

	return 0, nil
}
