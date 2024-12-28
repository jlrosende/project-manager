package repositories

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/config"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type ProjectRepository struct {
	git *config.Config
}

var _ ports.ProjectRepository = (*ProjectRepository)(nil)

func NewProjectRepository() (*ProjectRepository, error) {
	git, err := config.LoadConfig(config.GlobalScope)

	if err != nil {
		return nil, err
	}

	return &ProjectRepository{
		git: git,
	}, nil
}

func (p *ProjectRepository) Get(name string) (*domain.Project, error) {

	return nil, nil
}

func (p *ProjectRepository) List() ([]domain.Project, error) {

	projects := []domain.Project{}
	for _, section := range p.git.Raw.Sections {
		if section.IsName("includeIf") {
			for _, sub := range section.Subsections {
				if sub.HasOption("project") {
					projects = append(projects, domain.Project{
						Name: sub.Option("project"),
						Path: filepath.Dir(sub.Option("path")),
					})
				}
			}

		}
	}

	return projects, nil

}

func (p *ProjectRepository) Create(name, path, subproject string, env_vars map[string]string) (*domain.Project, error) {
	// Get the global config

	path, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	log.Println(path)

	gitdir := fmt.Sprintf("gitdir/i:%s/", filepath.Join(path))

	// Create git this if not exist

	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Println(err)
	}

	gitConfigName := fmt.Sprintf(".%s.gitconfig", name)

	gitConfigPath := filepath.Join(path, gitConfigName)

	includeIf := p.git.Raw.Section("includeIf").Subsection(gitdir)

	includeIf.SetOption("path", gitConfigPath)
	includeIf.SetOption("project", name)

	if subproject != "" {
		includeIf.SetOption("subproject", subproject)
	}

	gitConf, _ := p.git.Marshal()

	log.Printf("\n%s", string(gitConf))

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filepath.Join(home, ".gitconfig")

	fpGit, err := os.OpenFile(filepath.Join(home, ".gitconfig"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return nil, err
	}
	defer fpGit.Close()

	n_bytes, err := fpGit.Write(gitConf)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("%d Bytes written in %s", n_bytes, gitConfigPath))

	if err := p.git.Validate(); err != nil {
		return nil, err
	}

	// New Git config

	newConfig := config.NewConfig()

	newConfig.User.Email = "test-3"
	newConfig.User.Name = "test"

	newConfig.Raw.AddOption("commit", "", "gpgsign", "true")

	newGitConf, err := newConfig.Marshal()

	if err != nil {
		return nil, err
	}

	_, err = os.Stat(gitConfigPath)

	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exists in directory %s", gitConfigName, path)
	}

	log.Println(string(gitConfigPath))
	log.Printf("\n%s", string(newGitConf))

	fp, err := os.OpenFile(gitConfigPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	n_bytes, err = fp.Write(newGitConf)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("%d Bytes written in %s", n_bytes, gitConfigPath))

	if err != nil {
		return nil, err
	}

	if err := newConfig.Validate(); err != nil {
		return nil, err
	}

	// New .env file with env_vars
	envPath := filepath.Join(path, ".env")
	_, err = os.Stat(envPath)

	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exists in directory %s", envPath, path)
	}

	fpEnv, err := os.OpenFile(envPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return nil, err
	}
	defer fpEnv.Close()

	for key, value := range env_vars {
		n_bytes_env, err := fpEnv.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return nil, err
		}
		slog.Info(fmt.Sprintf("%d Bytes written in %s", n_bytes_env, envPath))

	}

	// n_bytes_env, err := fpEnv.Write(envVars)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func (p *ProjectRepository) Edit(name string) (*domain.Project, error) {
	return nil, nil

}

func (p *ProjectRepository) Delete(name string) error {
	return nil

}
