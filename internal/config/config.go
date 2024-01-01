package config

import (
	"errors"
	"fmt"

	"sigs.k8s.io/yaml"
)

// ErrInvalidConfig configuration is invalid.
var ErrInvalidConfig = errors.New("invalid runtime config")

// New create a new instance of Config.
func New() *Config {
	return &Config{}
}

// GetConfig shallow load the resource.
func (c *Config) GetConfig(
	fileContent []byte,
) (*Config, error) {
	config := New()

	if err := yaml.Unmarshal(fileContent, config); err != nil {
		return nil, fmt.Errorf("%w: cannot unmarshal config", ErrInvalidConfig)
	}

	return config, nil
}
