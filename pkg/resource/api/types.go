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

// Manager interface all resources must implement.
type Manager interface {
	Reconcile() error
	GetStatus() Phase
	GetMessage() string
	GetReason() string
}
