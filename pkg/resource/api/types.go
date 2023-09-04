package api

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
)

// Phase is a label for the condition of the resource at the current time.
type Phase string

// Status contains status of the resource.
type Status struct {
	// Phase sets `phase` as .status.Phase of the resource.
	Phase Phase
	// A human readable message indicating details about the transition.
	Message string
	// The reason for the condition's last transition.
	Reason string
	// Conditions contains status of the File lifecycle.
	Conditions []StatusConditions
}

// StatusConditions contains status of the resource lifecycle.
type StatusConditions struct {
	// Type the resources condition type.
	Type string
	// Status the resources phase.
	Status Phase
	// A human readable message indicating details about the transition.
	Message string
	// The reason for the condition's last transition.
	Reason string
	// Got the resources current state.
	Got string
	// Want the resources desired state.
	Want string
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
func (sc *StatusConditions) GetReason() string { return sc.Reason }

// SetReason set the reason property.
func (sc *StatusConditions) SetReason(reason string) { sc.Reason = reason }

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
		reason string,
		got string,
		want string,
	)
}
