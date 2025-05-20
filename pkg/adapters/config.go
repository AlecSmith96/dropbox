package adapters

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
)

type Config struct {
	SourceDirectory      string `yaml:"source-directory"`
	DestinationDirectory string `yaml:"destination-directory"`
	BaseURL              string `yaml:"base-url"`
	UseAbsolutePaths     bool   `yaml:"use-absolute-paths"`
}

func NewConfig() (*Config, error) {
	var config Config
	err := cleanenv.ReadConfig("dev-config.yaml", &config)
	if err != nil {
		slog.Debug("no config file found, falling back to env vars", "err", err)
	}

	return &config, nil
}
