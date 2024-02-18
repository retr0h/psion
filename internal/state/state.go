package state

import (
	"encoding/json"
	"os"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/pkg/resource/api"
)

// New create a new instance of State.
func New(
	fileManager internal.FileManager,
	stateFile string,
) *State {
	return &State{
		FileName:    stateFile,
		fileManager: fileManager,
	}
}

// GetStatus determine the state status.
func (s *State) GetStatus() api.Phase {
	noop := s.allMatch(api.NoOp)
	// set status to `NoOp` when all status are `NoOp`
	if noop {
		return api.NoOp
	}

	succeeded := s.allMatch(api.Succeeded)
	// set status to `Succeeded` when all status are `Succeeded`
	if succeeded {
		return api.Succeeded
	}

	pending := s.anyMatch(api.Pending)
	// set status to `Pending` when any status `Pending`
	if pending {
		return api.Pending
	}

	failed := s.anyMatch(api.Failed)
	// set status to `Failed` when any status are `Failed`
	if failed {
		return api.Failed
	}

	// otherwise set to `Unknown`
	return api.Unknown
}

// GetItems get the items property.
func (s *State) GetItems() []*StateResource { return s.Items }

// SetItems set the items property.
func (s *State) SetItems(
	stateResource *StateResource,
) {
	stateResources := s.Items
	stateResources = append(stateResources, stateResource)
	s.Items = stateResources
}

// GetStatusString the status property as a string.
func (s *State) GetStatusString() string { return string(s.GetStatus()) }

func (s *State) allMatch(phase api.Phase) bool {
	for _, resource := range s.Items {
		if resource.GetStatus() != phase {
			return false
		}
	}

	return true
}

func (s *State) anyMatch(phase api.Phase) bool {
	for _, resource := range s.Items {
		if resource.GetStatus() == phase {
			return true
		}
	}
	return false
}

// SetState write state to file.
func (s *State) SetState() error {
	state := State{
		Items: s.Items,
	}
	b, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(s.FileName, b, 0o644)
}

// GetState read state from file and unmarshal.
func (s *State) GetState() (*State, error) {
	fileContent, err := s.fileManager.Read(s.FileName)
	if err != nil {
		return nil, err
	}

	state := &State{}
	if err := json.Unmarshal(fileContent, state); err != nil {
		return nil, err
	}

	return state, nil
}

// GetStatus the status property.
func (r *StateResource) GetStatus() api.Phase { return r.Status.Phase }
