package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %v", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)

	return &cfg, nil
}
