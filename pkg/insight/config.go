package insight

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Config holds the configuration of multiple projects for which metrics need to be collected.
type Config struct {
	Sources  []Source `json:"sources,omitempty" yaml:"sources,omitempty"`
	LogLevel string   `json:"log_level,omitempty" yaml:"log_level,omitempty"`
	logger   *logrus.Logger
}

// GetConfig returns the new instance of Config.
func GetConfig(conf Config, path string) (*Config, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Printf("config file %s not found, dropping configurations from file", path)

		return &conf, fmt.Errorf("fetching config file information failed with: %w", err)
	}

	fileOUT, err := os.ReadFile(path)
	if err != nil {
		log.Println("failed to read the config file, dropping configurations from file")

		return &conf, fmt.Errorf("reading config file errored with: %w", err)
	}

	var newConfig Config
	if err = yaml.Unmarshal(fileOUT, &newConfig); err != nil {
		log.Println("failed to unmarshall configurations, dropping configurations from file")

		return &conf, fmt.Errorf("parsing config file errored with: %w", err)
	}
	if err = mergo.Merge(&conf, &newConfig); err != nil {
		log.Println("failed to merge configurations, dropping configurations from file")

		return &conf, fmt.Errorf("merging config errored with: %w", err)
	}

	return &newConfig, nil
}

func (cfg *Config) SetLogger(logger *logrus.Logger) {
	cfg.logger = logger
}
