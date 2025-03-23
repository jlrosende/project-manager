package repositories

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/hcl/v2/hclsimple"
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

	for _, section := range p.git.Raw.Sections {
		if section.IsName("includeIf") {
			slog.Debug(fmt.Sprintf("Section: %+v\n", section))
			for _, sub := range section.Subsections {
				// Read subsections and get path of the project
				slog.Debug(fmt.Sprintf("\t - SubSection: %+v\n", sub))

				if path, ok := strings.CutPrefix(sub.Name, "gitdir/i:"); ok {
					// read in path .project
					projetPath := filepath.Join(path, ".project.hcl")

					project, err := loadDotProject(projetPath)

					if err != nil {
						return nil, err
					}

					if project.Name != name {
						continue
					}

					project.Path = path

					return project, nil
				}

			}
		}
	}

	return nil, fmt.Errorf("project '%s' not found, projects are case sensitive", name)
}

func (p *ProjectRepository) List() ([]*domain.Project, error) {

	projects := []*domain.Project{}
	for _, section := range p.git.Raw.Sections {
		if section.IsName("includeIf") {
			for _, sub := range section.Subsections {

				if path, ok := strings.CutPrefix(sub.Name, "gitdir/i:"); ok {

					projetPath := filepath.Join(path, ".project.hcl")

					project, err := loadDotProject(projetPath)
					if err != nil {
						slog.Error(err.Error())
						continue
					}

					project.Path = path

					projects = append(projects, project)
				}

			}

		}
	}

	return projects, nil

}

/*
TODO
  - Create directory if not exist
  - Create .env .project.hcl and .<project>.gitconfig files
  - If .env exist warn and continue
  - If .<project>.gitconfig exist warn and continue
  - If .project.hcl exist
*/
func (p *ProjectRepository) Create(name, path, subproject string, envVars domain.EnvVars, gitConfig *domain.GitConfig) (*domain.Project, error) {

	// Check if the path is a directory

	if info, err := os.Stat(path); err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, errors.New("the path must be a directory")
	}

	path, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	log.Println(path)

	gitdir := fmt.Sprintf("gitdir/i:%s/", filepath.Join(path))

	// Create paths this if not exist
	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Println(err)
	}

	// Check if the path is empty

	if isEmpty, err := IsDirEmpty(path); err != nil {
		return nil, err
	} else if !isEmpty {
		return nil, fmt.Errorf("directory %s, is not empty", path)
	}

	gitConfigName := fmt.Sprintf(".%s.gitconfig", name)

	gitConfigPath := filepath.Join(path, gitConfigName)

	includeIf := p.git.Raw.Section("includeIf").Subsection(gitdir)

	includeIf.SetOption("path", gitConfigPath)
	// TODO the project name is stored in a .project file
	// includeIf.SetOption("project", name)

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

	// TODO Check if values are correct or set
	newConfig.User.Email = gitConfig.User.Email
	newConfig.User.Name = gitConfig.User.Name

	newConfig.Raw.AddOption("user", "", "signingkey", gitConfig.User.SigningKey)

	newConfig.Raw.AddOption("commit", "", "gpgsign", strconv.FormatBool(gitConfig.Commit.GPGSign))
	newConfig.Raw.AddOption("tag", "", "gpgsign", strconv.FormatBool(gitConfig.Tag.GPGSign))

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

	for key, value := range envVars {
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

func (p *ProjectRepository) Delete(name string) error {
	return nil
}

func loadDotProject(path string) (*domain.Project, error) {
	project := &domain.Project{}

	if strings.HasPrefix(path, "~/") {
		dirname, _ := os.UserHomeDir()
		path = filepath.Join(dirname, path[2:])
	}

	err := hclsimple.DecodeFile(path, nil, project)
	if err != nil {
		return nil, err
	}

	if project.Shell == "" {
		if shell, ok := os.LookupEnv("SHELL"); ok {
			project.Shell = shell
		} else {
			return nil, fmt.Errorf("unable to load default shell. Add SHELL environment variable or set the `shell` variable inside your %s file", filepath.Join(path, ".project.hcl"))
		}
	}

	return project, nil
}

func IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
