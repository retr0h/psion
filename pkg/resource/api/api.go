package api

// Phase is a label for the condition of the resource at the current time.
type Phase string

// Action the kind of action to perform on resources.
type Action string

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
	Type string `json:"type,omitempty"`
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

// GetType the type property.
func (sc *StatusConditions) GetType() string { return sc.Type }

// GetStatus the status property.
func (sc *StatusConditions) GetStatus() Phase { return sc.Status }

// GetMessage get the message property.
func (sc *StatusConditions) GetMessage() string { return sc.Message }

// SetMessage set the message property.
func (sc *StatusConditions) SetMessage(message string) { sc.Message = message }

// GetReason the reason property.
func (sc *StatusConditions) GetReason() Action { return sc.Reason }

// GetReasonString the reason property.
func (sc *StatusConditions) GetReasonString() string { return string(sc.Reason) }

// GetStatusString the status property as a string.
func (sc *StatusConditions) GetStatusString() string { return string(sc.Status) }

// GetGot get the got property.
func (sc *StatusConditions) GetGot() string { return sc.Got }

// GetWant get the want property.
func (sc *StatusConditions) GetWant() string { return sc.Want }

// Manager interface all resources must implement.
type Manager interface {
	Reconcile() error
	GetStatus() Phase
	GetStatusString() string
	SetStatus(status Phase)
	GetStatusConditions() []StatusConditions
	SetStatusCondition(
		statusType string,
		status Phase,
		message string,
		got string,
		want string,
	)
	GetState() *Resource
}

// State used by state file and status.
type State struct {
	Items []*Resource `json:"items,omitempty"`
}

// GetStatus determine the state status.
func (s *State) GetStatus() Phase {
	noop := s.allStatusMatch(NoOp, s.Items)
	// set status to `NoOp` when all status are `NoOp`
	if noop {
		return NoOp
	}

	succeeded := s.allStatusMatch(Succeeded, s.Items)
	// set status to `Succeeded` when all status are `Succeeded`
	if succeeded {
		return Succeeded
	}

	pending := s.anyStatusMatch(Pending, s.Items)
	// set status to `Pending` when any status `Pending`
	if pending {
		return Pending
	}

	failed := s.anyStatusMatch(Failed, s.Items)
	// set status to `Failed` when any status are `Failed`
	if failed {
		return Failed
	}

	// otherwise set to `Unknown`
	return Unknown
}

// GetStatusString the status property as a string.
func (s *State) GetStatusString() string { return string(s.GetStatus()) }

func (s *State) allStatusMatch(phase Phase, resources []*Resource) bool {
	for _, resource := range resources {
		if resource.GetStatus() != phase {
			return false
		}
	}

	return true
}

func (s *State) anyStatusMatch(phase Phase, resources []*Resource) bool {
	for _, resource := range resources {
		if resource.GetStatus() == phase {
			return true
		}
	}
	return false
}

// Resource container holding resource state.
type Resource struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`

	Phase  Phase   `json:"phase,omitempty"`
	Status *Status `json:"status"`
}

// GetStatus the status property.
func (r *Resource) GetStatus() Phase { return r.Status.Phase }
