package api

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
