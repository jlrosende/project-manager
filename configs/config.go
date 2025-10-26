package configs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Theme           string        `hcl:"theme"`
	BackupDirectory string        `hcl:"backup_directory,optional"`
	CustomThemes    []CustomTheme `hcl:"custom_theme,block"`
}

type CustomTheme struct {
	Name        string `hcl:"name,label"`
	Title       string `hcl:"title,optional"`
	Section     string `hcl:"section,optional"`
	Subtext     string `hcl:"subtext,optional"`
	Text        string `hcl:"text,optional"`
	Placeholder string `hcl:"placeholder,optional"`
	Border      string `hcl:"border,optional"`
	Error       string `hcl:"error,optional"`
	ButtonDefFg string `hcl:"buttonDefFg,optional"`
	ButtonDefBg string `hcl:"buttonDefBg,optional"`
	ButtonSelFg string `hcl:"buttonSelFg,optional"`
	ButtonSelBg string `hcl:"buttonSelBg,optional"`
	SelectedFg  string `hcl:"selectedFg,optional"`
	SelectedBg  string `hcl:"selectedBg,optional"`
	Help        string `hcl:"help,optional"`
}

func GetConfig(cfgFile string) (*Config, error) {
	if cfgFile == "" {
		if env := os.Getenv("PM_CONFIG"); strings.TrimSpace(env) != "" {
			cfgFile = env
		}
	}

	if strings.HasPrefix(cfgFile, "~/") {
		h, _ := os.UserHomeDir()
		cfgFile = filepath.Join(h, cfgFile[2:])
	}

	if strings.TrimSpace(cfgFile) != "" {
		var cfg Config
		if err := hclsimple.DecodeFile(cfgFile, nil, &cfg); err != nil {
			return nil, err
		}

		applyDefaults(&cfg)

		return &cfg, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		cfg := &Config{}
		applyDefaults(cfg)

		return cfg, nil
	}

	defaultPath := filepath.Join(configDir, "pm", "config.hcl")
	if _, err := os.Stat(defaultPath); err == nil {
		var cfg Config
		if err := hclsimple.DecodeFile(defaultPath, nil, &cfg); err != nil {
			return nil, err
		}

		applyDefaults(&cfg)

		return &cfg, nil
	}

	cfg := &Config{}
	applyDefaults(cfg)

	return cfg, nil
}

func applyDefaults(cfg *Config) {
	if strings.TrimSpace(cfg.Theme) == "" {
		cfg.Theme = "nord"
	}

	if strings.TrimSpace(cfg.BackupDirectory) == "" {
		cfg.BackupDirectory = filepath.Join("~", ".pm", "backups")
	}
}
