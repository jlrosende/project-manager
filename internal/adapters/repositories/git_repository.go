package repositories

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/config"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type GitRepository struct {
	git    *config.Config
	global *config.Config
}

var _ ports.GitRepository = (*GitRepository)(nil)

func NewGitRepository() (*GitRepository, error) {
	return &GitRepository{
		git: config.NewConfig(),
	}, nil
}

func (g *GitRepository) Load(path string) (*domain.GitConfig, error) {
	gitFileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := g.git.Unmarshal(gitFileContent); err != nil {
		return nil, err
	}

	if err := g.git.Validate(); err != nil {
		return nil, err
	}

	tagOpt := g.git.Raw.Section("tag").Option("gpgsign")
	tagGPGSing := false

	if tagOpt != "" {
		if b, e := strconv.ParseBool(tagOpt); e == nil {
			tagGPGSing = b
		}
	}

	commitOpt := g.git.Raw.Section("commit").Option("gpgsign")
	commitGPGSing := false

	if commitOpt != "" {
		if b, e := strconv.ParseBool(commitOpt); e == nil {
			commitGPGSing = b
		}
	}

	return &domain.GitConfig{
		User: domain.User{
			Name:       g.git.User.Name,
			Email:      g.git.User.Email,
			SigningKey: g.git.Raw.Section("user").Option("signingkey"),
		},
		Tag: domain.Tag{
			GPGSign: tagGPGSing,
		},
		Commit: domain.Commit{
			GPGSign: commitGPGSing,
		},
	}, nil
}

func (g *GitRepository) Save(path string, gitConfig *domain.GitConfig) error {
	g.git.User.Email = gitConfig.User.Email
	g.git.User.Name = gitConfig.User.Name
	g.git.Raw.AddOption("user", "", "signingkey", gitConfig.User.SigningKey)

	g.git.Raw.AddOption("tag", "", "gpgsign", "true")
	g.git.Raw.AddOption("commit", "", "gpgsign", "true")

	newGitConf, err := g.git.Marshal()
	if err != nil {
		return err
	}

	_, _ = os.Stat(path)

	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(newGitConf)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitRepository) LoadGlobal() error {
	c, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		return err
	}

	g.global = c

	return nil
}

func (g *GitRepository) UpdateIncludeIf(gitdir, perProjectPath, subproject string) error {
	if g.global == nil {
		if err := g.LoadGlobal(); err != nil {
			return err
		}
	}

	includeIf := g.global.Raw.Section("includeIf").Subsection(gitdir)
	includeIf.SetOption("path", perProjectPath)

	if subproject != "" {
		includeIf.SetOption("subproject", subproject)
	}

	return nil
}

func (g *GitRepository) SaveGlobal(home string) error {
	if g.global == nil {
		if err := g.LoadGlobal(); err != nil {
			return err
		}
	}

	b, err := g.global.Marshal()
	if err != nil {
		return err
	}

	fp, err := os.OpenFile(filepath.Join(home, ".gitconfig"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err := fp.Write(b); err != nil {
		return err
	}

	return g.global.Validate()
}

func (g *GitRepository) RemoveHooks(_ context.Context, project domain.ProjectIdentifier) error {
	if strings.TrimSpace(project.Path) == "" {
		return nil
	}

	name := strings.TrimSpace(project.Name)
	if name == "" {
		name = filepath.Base(project.Path)
	}

	perProject := filepath.Join(project.Path, fmt.Sprintf(".%s.gitconfig", name))
	if err := os.Remove(perProject); err != nil && !os.IsNotExist(err) {
		return err
	}

	hooksDir := filepath.Join(project.Path, ".git", "hooks")
	entries, err := os.ReadDir(hooksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	for _, entry := range entries {
		full := filepath.Join(hooksDir, entry.Name())
		if entry.IsDir() {
			if err := os.RemoveAll(full); err != nil {
				return err
			}
			continue
		}

		if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}
