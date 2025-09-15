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

					if strings.HasPrefix(path, "~/") {
						dirname, _ := os.UserHomeDir()
						path = filepath.Join(dirname, path[2:])
					}

					project.Path = path

					slog.Debug("load project", "path", project.Path)

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
	// Expand '~' to user home
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}

	// Ensure path exists and is a directory (create if missing)
	if info, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else if !info.IsDir() {
		return nil, errors.New("the path must be a directory")
	}

	// Normalize to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	path = absPath

	log.Println(path)

	gitdir := fmt.Sprintf("gitdir/i:%s/", filepath.Join(path))

	// Create paths (idempotent)
	if err := os.MkdirAll(path, 0o755); err != nil {
		return nil, err
	}

	// Require empty directory to avoid clobbering existing projects
	if isEmpty, err := IsDirEmpty(path); err != nil {
		return nil, err
	} else if !isEmpty {
		return nil, fmt.Errorf("directory %s, is not empty", path)
	}

	gitConfigName := fmt.Sprintf(".%s.gitconfig", name)
	gitConfigPath := filepath.Join(path, gitConfigName)

	includeIf := p.git.Raw.Section("includeIf").Subsection(gitdir)
	includeIf.SetOption("path", gitConfigPath)
	if subproject != "" {
		includeIf.SetOption("subproject", subproject)
	}

	gitConf, _ := p.git.Marshal()
	log.Printf("\n%s", string(gitConf))

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fpGit, err := os.OpenFile(filepath.Join(home, ".gitconfig"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	defer fpGit.Close()
	if _, err := fpGit.Write(gitConf); err != nil {
		return nil, err
	}
	if err := p.git.Validate(); err != nil {
		return nil, err
	}

	// Per-project git config
	newConfig := config.NewConfig()
	newConfig.User.Email = gitConfig.User.Email
	newConfig.User.Name = gitConfig.User.Name
	newConfig.Raw.AddOption("user", "", "signingkey", gitConfig.User.SigningKey)
	newConfig.Raw.AddOption("commit", "", "gpgsign", strconv.FormatBool(gitConfig.Commit.GPGSign))
	newConfig.Raw.AddOption("tag", "", "gpgsign", strconv.FormatBool(gitConfig.Tag.GPGSign))
	newGitConf, err := newConfig.Marshal()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(gitConfigPath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exists in directory %s", gitConfigName, path)
	}
	fp, err := os.OpenFile(gitConfigPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	if _, err := fp.Write(newGitConf); err != nil {
		return nil, err
	}
	if err := newConfig.Validate(); err != nil {
		return nil, err
	}

	// .env file with env_vars
	envPath := filepath.Join(path, ".env")
	if _, err = os.Stat(envPath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exists in directory %s", envPath, path)
	}
	fpEnv, err := os.OpenFile(envPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	defer fpEnv.Close()
	for key, value := range envVars {
		if _, err := fmt.Fprintf(fpEnv, "%s=%s\n", key, value); err != nil {
			return nil, err
		}
	}

	// .project.hcl
	projPath := filepath.Join(path, ".project.hcl")
	if _, err = os.Stat(projPath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exists in directory %s", filepath.Base(projPath), path)
	}
	fpProj, err := os.OpenFile(projPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	defer fpProj.Close()
	projectHCL := fmt.Sprintf("name = \"%s\"\n\ndescription = \"\"\n\nenv_vars_file = \".env\"\n", name)
	if _, err := fpProj.WriteString(projectHCL); err != nil {
		return nil, err
	}

	return &domain.Project{
		Name:        name,
		Description: "",
		Path:        path,
		EnvVarsFile: ".env",
	}, nil
}

func (p *ProjectRepository) Delete(name string) error {
	return nil
}

func (p *ProjectRepository) AddEnvironment(projectName string, env *domain.Environment, envVars domain.EnvVars) error {
	project, err := p.Get(projectName)
	if err != nil {
		return err
	}
	for _, e := range project.Environments {
		if e.Name == env.Name {
			return fmt.Errorf("environment %s already exists", env.Name)
		}
	}
	if env.EnvVarsMode == "" {
		env.EnvVarsMode = domain.ENV_VARS_MODE_MERGE
	}
	if strings.HasPrefix(project.Path, "~/") {
		h, _ := os.UserHomeDir()
		project.Path = filepath.Join(h, project.Path[2:])
	}
	envFile := env.EnvVarsFile
	if envFile == "" {
		envFile = ".env." + env.Name
	}
	envPath := envFile
	if !filepath.IsAbs(envFile) {
		envPath = filepath.Join(project.Path, envFile)
	}
	if _, err := os.Stat(envPath); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists in directory %s", filepath.Base(envPath), filepath.Dir(envPath))
	}
	fpEnv, err := os.OpenFile(envPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer fpEnv.Close()
	for k, v := range envVars {
		if _, err := fmt.Fprintf(fpEnv, "%s=%s\n", k, v); err != nil {
			return err
		}
	}
	projPath := filepath.Join(project.Path, ".project.hcl")
	fp, err := os.OpenFile(projPath, os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer fp.Close()
	b := &strings.Builder{}
	b.WriteString("\n")
	b.WriteString("environment \"")
	b.WriteString(env.Name)
	b.WriteString("\" {\n")
	if env.Color != "" {
		b.WriteString("  color = \"")
		b.WriteString(env.Color)
		b.WriteString("\"\n")
	}
	b.WriteString("  env_vars_mode = \"")
	b.WriteString(env.EnvVarsMode)
	b.WriteString("\"\n")
	b.WriteString("  env_vars_file = \"")
	b.WriteString(envFile)
	b.WriteString("\"\n")
	b.WriteString("}\n")
	if _, err := io.WriteString(fp, b.String()); err != nil {
		return err
	}
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
