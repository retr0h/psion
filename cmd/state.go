package cmd

import (
	"github.com/retr0h/psion/internal/state"
	"github.com/retr0h/psion/pkg/resource/api"
)

// StateManager interface to state.
type StateManager interface {
	GetStatus() api.Phase
	GetStatusString() string
	GetItems() []*state.StateResource
	SetItems(stateResource *state.StateResource)
	SetState() error
	GetState() (*state.State, error)
}
