package api

import (
	"github.com/retr0h/psion/internal"
)

// Phase is a label for the condition of the resource at the current time.
type Phase string

// Action plan or apply the resources.
type Action string

// SpecAction the kind of action to perform on resources.
type SpecAction string

// These are the valid statuses of the resource.
const (
	// Pending means the declared changes have yet to be made.
	Pending Phase = "Pending"
	// Succeeded means the declared changes have been made.
	Succeeded Phase = "Succeeded"
	// Failed means the declaared changes have not been made.
	Failed Phase = "Failed"
	// Unknown means that for some reason the state of the resource
	// could not be obtained.
	Unknown Phase = "Unknown"
	// NoOp means no changes will be made, resource matches declared state.
	NoOp Phase = "NoOp"

	// Plan represents the changes to make consistent with the desired state.
	Plan Action = "Plan"
	// Apply represents the changes to make the desired state.
	Apply Action = "Apply"
)

// Status contains status of the resource.
type Status struct {
	// Phase sets `phase` as .status.Phase of the resource.
	Phase Phase `json:"phase,omitempty"`
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
	Type SpecAction `json:"type,omitempty"`
	// Status the resources phase.
	Status Phase `json:"status,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
	// The reason for the condition's last transition.
	Reason Action `json:"reason,omitempty"`
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

	Phase  Phase   `json:"phase,omitempty"`
	Status *Status `json:"status"`
}
