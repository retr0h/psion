package config

import (
	"embed"
	"fmt"

	"sigs.k8s.io/yaml"

	"github.com/retr0h/psion/internal"
)

// New create a new instance of Config.
func New(
	eFs embed.FS,
	fileManager internal.FileManager,
) *Config {
	return &Config{
		fm:  fileManager,
		eFs: eFs,
	}
}

// LoadConfig shallow load the resource.
func (c *Config) LoadConfig(
	fileContent []byte,
) error {
	if err := yaml.Unmarshal(fileContent, c); err != nil {
		return fmt.Errorf("%w: cannot unmarshal config", err)
	}

	return nil
}

// GetName the Name property.
func (c *Config) GetName() string { return c.Name }

// GetAPIVersion the APIVersion property.
func (c *Config) GetAPIVersion() string { return c.APIVersion }

// GetKind the Kind property.
func (c *Config) GetKind() string { return c.Kind }
