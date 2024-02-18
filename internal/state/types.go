package state

import (
	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/status"
	"github.com/retr0h/psion/pkg/resource/api"
)

// State used by state file and status.
type State struct {
	Items       []*StateResource     `json:"items,omitempty"`
	FileName    string               `json:"-"`
	fileManager internal.FileManager `json:"-"`
}

// StateResource container holding resource state.
type StateResource struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`

	Phase  api.Phase      `json:"phase,omitempty"`
	Status *status.Status `json:"status"`
}
