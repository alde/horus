package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// Config struct
type Config struct {
	Horus struct {
		Host string
	}

	GoogleKMS struct {
		Project  string
		KeyName  string
		KeyRing  string
		Location string
	}
}

// New creates a new Config
func New(configFile string) *Config {
	cfg := &Config{}
	ReadConfigFile(cfg, getConfigFilePath(configFile))

	return cfg
}

// getConfigFilePath returns the location of the config file in order of priority:
// 1 ) --config commandline flag
// 1 ) $(HOME)/.config/horus/cli_config.toml
func getConfigFilePath(configPath string) string {
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
		panic(fmt.Sprintf("Unable to open %s.", configPath))
	}
	path, _ := os.UserHomeDir()
	path = fmt.Sprintf("%s/.config/horus/cli_config.toml", path)
	if _, err := os.Open(path); err == nil {
		return path
	}

	return ""
}

// ReadConfigFile reads the config file and merges with DefaultConfig, taking precedence
func ReadConfigFile(cfg *Config, path string) {
	_, err := os.Stat(path)
	if err != nil {
		return
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		log.Fatal("unable to read config")
	}
}
