package configs

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// go:embed config.default.hcl
var Default []byte

type Config struct {
	Theme      string `mapstructure:"theme"`
	RootFolder string `mapstructure:"root_folder"`
}

func NewConfig(cfgFile string) *Config {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).

		viper.AddConfigPath(path.Join(home, ".config/pm/"))
		viper.SetConfigName("config")
		viper.SetConfigType("hcl")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	config := new(Config)

	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatal("Failed to unmarshal config", l.Error(err))
	}

	return config
}
