package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

const (
	defaultConfigFilePath     string = "config.yaml"
	configFilePathEnvVariable string = "CONFIG_FILEPATH"
)

func Load() (*Config, error) {
	configFilePath := getConfigFilePath()
	data, err := os.ReadFile(getConfigFilePath())
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %v", configFilePath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)

	return &cfg, nil
}

func getConfigFilePath() string {
	configFilePath, isEnvSet := os.LookupEnv(configFilePathEnvVariable)
	if isEnvSet {
		return configFilePath
	}
	return configFilePathEnvVariable
}
