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

	"github.com/spf13/viper"
)

// go:embed config.default.hcl
var defaultConfig []byte

var (
	NotFoundErr = errors.New("config file not found")
)

type Config struct {
	Theme      string `mapstructure:"theme"`
	RootFolder string `mapstructure:"root_folder"`
}

func GetConfig(cfgFile string) (*Config, error) {

	v, err := LoadConfig(cfgFile)

	if err != nil {
		if errors.Is(err, NotFoundErr) {
			defaultViper, err := DefaultConfig()

			if err != nil {
				return nil, err
			}

			v = defaultViper
		}
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
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config

	err := v.Unmarshal(&cfg)

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
