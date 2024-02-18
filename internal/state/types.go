package state

import (
	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/pkg/resource/api"
)

// Status contains status of the resource.
type Status struct {
	// Phase sets `phase` as .status.Phase of the resource.
	Phase api.Phase `json:"phase,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Conditions contains status of the File lifecycle.
	Conditions []StatusConditions `json:"conditions,omitempty"`
}

// StatusConditions contains status of the resource lifecycle.
type StatusConditions struct {
	// Type the resources condition type.
	Type api.SpecAction `json:"type,omitempty"`
	// Status the resources phase.
	Status api.Phase `json:"status,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
	// The reason for the condition's last transition.
	Reason api.Action `json:"reason,omitempty"`
	// Got the resources current state.
	Got string `json:"got,omitempty"`
	// Want the resources desired state.
	Want string `json:"want,omitempty"`
}

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

	Phase  api.Phase `json:"phase,omitempty"`
	Status *Status   `json:"status"`
}
