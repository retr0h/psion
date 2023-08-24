package config

import (
	"errors"
	"fmt"

	"sigs.k8s.io/yaml"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ErrInvalidRuntimeConfig = errors.New("invalid runtime config")

// Runtime config is a static copy of the resource.
type RuntimeConfig struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	metav1.TypeMeta   `json:",omitempty,inline"`
}

// LoadRuntimeConfig shallow load the resource.
func LoadRuntimeConfig(fileContent []byte) (*RuntimeConfig, error) {
	config := &RuntimeConfig{}

	if err := yaml.Unmarshal(fileContent, config); err != nil {
		return nil, fmt.Errorf("%w: cannot unmarshal config", ErrInvalidRuntimeConfig)
	}

	return config, nil
}
