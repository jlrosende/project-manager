package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

var allowedProjectConfigKeys = map[string]struct{}{
	"name":         {},
	"here":         {},
	"path":         {},
	"environments": {},
	"metadata":     {},
}

// SkeletonFormat represents the serialization format for generated CLI input files.
type SkeletonFormat string

const (
	SkeletonFormatJSON SkeletonFormat = "json"
	SkeletonFormatYAML SkeletonFormat = "yaml"
)

// ProjectConfigFlags captures CLI-provided overrides for project creation.
type ProjectConfigFlags struct {
	Name    string
	Path    string
	Here    bool
	PathSet bool
	HereSet bool
}

// MergeProjectInputs combines configuration file inputs with CLI arguments,
// preferring explicit flag values. It returns the merged project definition and
// the environment variables destined for the .env file.
func MergeProjectInputs(
	cfg *domain.ConfigInput,
	flags ProjectConfigFlags,
) (domain.ProjectDefinition, domain.EnvVars, error) {
	var (
		def     domain.ProjectDefinition
		envVars = domain.EnvVars{}
	)

	if cfg != nil {
		if cfg.Name != nil {
			def.Name = strings.TrimSpace(*cfg.Name)
		}

		if cfg.Here != nil {
			def.Here = *cfg.Here
		}

		if cfg.Path != nil {
			def.Path = strings.TrimSpace(*cfg.Path)
		}

		if len(cfg.Environments) > 0 {
			def.Environments = copyStringMap(cfg.Environments)
			envVars = domain.EnvVars(copyStringMap(cfg.Environments))
		}

		if len(cfg.Metadata) > 0 {
			def.Metadata = copyStringMap(cfg.Metadata)
		}
	}

	if flags.PathSet {
		def.Path = strings.TrimSpace(flags.Path)
		if def.Here {
			def.Here = false
		}
	}

	if flags.HereSet {
		def.Here = flags.Here
		if flags.Here && !flags.PathSet {
			def.Path = ""
		}
	}

	if strings.TrimSpace(flags.Name) != "" {
		def.Name = strings.TrimSpace(flags.Name)
	}

	return def, envVars, nil
}

// LoadProjectConfig reads a JSON or YAML configuration file from disk.
func LoadProjectConfig(path string) (*domain.ConfigInput, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("config path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("config file %s is empty", path)
	}

	ext := strings.ToLower(filepath.Ext(path))

	var (
		cfg domain.ConfigInput
		raw map[string]any
	)

	switch ext {
	case ".json":
		raw, err = decodeProjectJSON(data, &cfg)
	case ".yaml", ".yml":
		raw, err = decodeProjectYAML(data, &cfg)
	default:
		raw, err = decodeProjectJSON(data, &cfg)
		if err != nil {
			raw, err = decodeProjectYAML(data, &cfg)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("parse config file %s: %w", path, err)
	}

	cfg.SetUnknownFields(filterProjectUnknown(raw))

	return &cfg, nil
}

func decodeProjectJSON(data []byte, cfg *domain.ConfigInput) (map[string]any, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var parsed domain.ConfigInput
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	*cfg = parsed

	return raw, nil
}

func decodeProjectYAML(data []byte, cfg *domain.ConfigInput) (map[string]any, error) {
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var parsed domain.ConfigInput
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	*cfg = parsed

	return raw, nil
}

func filterProjectUnknown(raw map[string]any) map[string]any {
	if len(raw) == 0 {
		return nil
	}

	unknown := map[string]any{}

	for k, v := range raw {
		if _, ok := allowedProjectConfigKeys[k]; !ok {
			unknown[k] = v
		}
	}

	if len(unknown) == 0 {
		return nil
	}

	return unknown
}

func copyStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}

	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}

	return out
}

// GenerateProjectSkeleton writes a CLI input skeleton in the requested format to the
// provided path. The generated file can be edited and passed to --cli-input for
// future project creation runs.
func GenerateProjectSkeleton(path string, format SkeletonFormat) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("skeleton output path is empty")
	}

	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return fmt.Errorf("%s is a directory", path)
		}

		return fmt.Errorf("file %s already exists", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check skeleton destination: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("prepare skeleton destination: %w", err)
	}

	data, err := RenderProjectSkeleton(format)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write skeleton file: %w", err)
	}

	return nil
}

// RenderProjectSkeleton returns the serialized CLI input skeleton for the requested format.
func RenderProjectSkeleton(format SkeletonFormat) ([]byte, error) {
	template := defaultProjectConfigSkeleton()

	var (
		data []byte
		err  error
	)

	switch format {
	case SkeletonFormatJSON:
		data, err = json.MarshalIndent(template, "", "  ")
	case SkeletonFormatYAML:
		data, err = yaml.Marshal(template)
	default:
		return nil, fmt.Errorf("unsupported skeleton format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("render skeleton: %w", err)
	}

	if !bytes.HasSuffix(data, []byte("\n")) {
		data = append(data, '\n')
	}

	return data, nil
}

func defaultProjectConfigSkeleton() domain.ConfigInput {
	name := "your-project-name"
	path := "/absolute/path/to/your-project"
	here := false

	return domain.ConfigInput{
		Name: &name,
		Path: &path,
		Here: &here,
		Environments: map[string]string{
			"EXAMPLE_ENV_VAR": "value",
		},
		Metadata: map[string]string{
			"description": "Describe your project",
		},
	}
}
