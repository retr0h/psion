package state

import (
	"github.com/retr0h/psion/pkg/resource/api"
)

// Move this
// Move this
// Move this
// Move this
// Move this

// Manager interface to state.
type Manager interface {
	GetStatus() api.Phase
	GetStatusString() string
	GetItems() []*api.StateResource
	SetItems(stateResource *api.StateResource)
	SetState() error
	GetState() (*api.State, error)
}
