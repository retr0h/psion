package internal

import (
	"github.com/retr0h/psion/internal/config"
)

// ConfigManager manager responsible for Config operations.
type ConfigManager interface {
	GetConfig(
		fileContent []byte,
	) (*config.Config, error)
}
