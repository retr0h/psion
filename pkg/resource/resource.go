package resource

import (
	"github.com/retr0h/psion/pkg/resource/api"
)

// Manager interface to resources.
type Manager interface {
	Reconcile() error
	GetStatus() api.Phase
	GetStatusString() string
	GetStatusConditions() []api.StatusConditions
	SetStatusCondition(
		statusType api.SpecAction,
		status api.Phase,
		message string,
		got string,
		want string,
	)
	GetState() *api.StateResource
}
