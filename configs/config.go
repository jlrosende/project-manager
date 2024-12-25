package configs

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

//go:embed config.default.hcl
var defaultConfig []byte

var (
	NotFoundErr = errors.New("config file not found")
)

type Config struct {
	Theme      string               `mapstructure:"theme"`
	RootFolder string               `mapstructure:"root_folder"`
	Projects   map[string][]Project `mapstructure:"project"`
}

type Project struct {
	Path         string                 `mapstructure:"path"`
	Theme        string                 `mapstructure:"theme"`
	EnvVars      map[string]string      `mapstructure:"env_vars"`
	EnvVarsFile  string                 `mapstructure:"env_vars_file"`
	Environments map[string]Environment `mapstructure:"environment"`
}

type Environment struct {
	Theme       string            `mapstructure:"theme"`
	EnvVars     map[string]string `mapstructure:"env_vars"`
	EnvVarsFile string            `mapstructure:"env_vars_file"`
}

func GetConfig(cfgFile string) (*Config, error) {

	v, err := LoadConfig(cfgFile)

	if err != nil {
		return nil, err
	}

	config, err := ParseConfig(v)

	if err != nil {
		log.Printf("Unable to parse config: %v", err)
		return nil, err
	}

	return config, nil
}

func LoadConfig(cfgFile string) (*viper.Viper, error) {
	v := viper.New()

	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {

		home, err := os.UserHomeDir()

		if err != nil {
			return nil, err
		}

		v.AddConfigPath(path.Join(home, ".config/pm/"))
		v.SetConfigName("config")
		v.SetConfigType("hcl")
	}

	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		slog.Debug(fmt.Sprintf("Unable to read config: %v", err))
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return DefaultConfig()
		}
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config

	configOption := viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		sliceOfMapsToMapHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	))

	err := v.Unmarshal(&cfg, configOption)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse config: %w", err)
	}

	return &cfg, nil
}

func DefaultConfig() (*viper.Viper, error) {
	v := viper.New()

	err := v.ReadConfig(bytes.NewReader(defaultConfig))

	if err != nil {
		return nil, err
	}

	return v, nil
}

// sliceOfMapsToMapHookFunc merges a slice of maps to a map
func sliceOfMapsToMapHookFunc() mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() == reflect.Slice && from.Elem().Kind() == reflect.Map && (to.Kind() == reflect.Struct || to.Kind() == reflect.Map) {
			source, ok := data.([]map[string]interface{})
			if !ok {
				return data, nil
			}
			if len(source) == 0 {
				return data, nil
			}
			if len(source) == 1 {
				return source[0], nil
			}
			// flatten the slice into one map
			convert := make(map[string]interface{})
			for _, mapItem := range source {
				for key, value := range mapItem {
					convert[key] = value
				}
			}
			return convert, nil
		}
		return data, nil
	}
}
