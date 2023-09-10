package api

import (
	"encoding/json"
	"os"

	"github.com/spf13/afero"

	"github.com/retr0h/psion/internal/file"
)

// GetType the action property.
func (sc *StatusConditions) GetType() SpecAction { return sc.Type }

// GetTypeString the action property.
func (sc *StatusConditions) GetTypeString() string { return string(sc.GetType()) }

// GetStatus the status property.
func (sc *StatusConditions) GetStatus() Phase { return sc.Status }

// GetMessage get the message property.
func (sc *StatusConditions) GetMessage() string { return sc.Message }

// SetMessage set the message property.
func (sc *StatusConditions) SetMessage(message string) { sc.Message = message }

// GetReason the reason property.
func (sc *StatusConditions) GetReason() Action { return sc.Reason }

// GetReasonString the reason property.
func (sc *StatusConditions) GetReasonString() string { return string(sc.GetReason()) }

// GetStatusString the status property as a string.
func (sc *StatusConditions) GetStatusString() string { return string(sc.GetStatus()) }

// GetGot get the got property.
func (sc *StatusConditions) GetGot() string { return sc.Got }

// GetWant get the want property.
func (sc *StatusConditions) GetWant() string { return sc.Want }

// NewState create a new instance of State.
func NewState(
	appFs afero.Fs,
	stateFile string,
) *State {
	return &State{
		File:  stateFile,
		appFs: appFs,
	}
}

// GetStatus determine the state status.
func (s *State) GetStatus() Phase {
	noop := s.allMatch(NoOp)
	// set status to `NoOp` when all status are `NoOp`
	if noop {
		return NoOp
	}

	succeeded := s.allMatch(Succeeded)
	// set status to `Succeeded` when all status are `Succeeded`
	if succeeded {
		return Succeeded
	}

	pending := s.anyMatch(Pending)
	// set status to `Pending` when any status `Pending`
	if pending {
		return Pending
	}

	failed := s.anyMatch(Failed)
	// set status to `Failed` when any status are `Failed`
	if failed {
		return Failed
	}

	// otherwise set to `Unknown`
	return Unknown
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

func (s *State) allMatch(phase Phase) bool {
	for _, resource := range s.Items {
		if resource.GetStatus() != phase {
			return false
		}
	}

	return true
}

func (s *State) anyMatch(phase Phase) bool {
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

	return os.WriteFile(s.File, b, 0o644)
}

// GetState read state from file and unmarshal.
func (s *State) GetState() (*State, error) {
	fileContent, err := file.Read(s.appFs, s.File)
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
func (r *StateResource) GetStatus() Phase { return r.Status.Phase }
