package internal

import (
	// remove
	"github.com/retr0h/psion/pkg/api"
)

// ConfigManager manager responsible for Config operations.
type ConfigManager interface {
	LoadConfig(
		fileContent []byte,
	) error
	GetName() string
	GetAPIVersion() string
	GetKind() string

	LoadAllEmbeddedResourceFiles(
		plan bool,
	) ([]api.ResourceManager, error)
	GetAllEmbeddedFileNames() ([]string, error)
}
