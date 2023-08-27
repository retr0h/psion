package v1alpha1

import (
	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
)

const (
	FileDoStatus    StateType = "Unknown"
	FileDoExists    StateType = "Exists"
	FileDoNotExists StateType = "DoesNotExists"

	FilePlanStatus      EventType = "FilePlanStatus"
	FilePlanExists      EventType = "FilePlanExists"
	FilePlanDoNotExists EventType = "FilePlanDoNotExists"
)

func FilePlanRemoveFSM() *StateMachine {
	return &StateMachine{
		States: States{
			Default: State{
				Events: Events{
					FilePlanStatus: FileDoStatus,
				},
			},
			FileDoStatus: State{
				Action: &FilePlanStatusAction{},
				Events: Events{
					FilePlanExists:      FileDoExists,
					FilePlanDoNotExists: FileDoNotExists,
				},
			},
			FileDoExists: State{
				Action: &FilePlanExistsAction{},
				Events: Events{
					FilePlanExists: FileDoExists,
				},
			},
			FileDoNotExists: State{
				Action: &FilePlanNoExistsAction{},
				Events: Events{
					FilePlanDoNotExists: FileDoNotExists,
				},
			},
		},
	}
}

// FilePlanStatusAction initial state when entering the state machine when
// Plan is set and the File resource should be removed.
type FilePlanStatusAction struct{}

// Execute implementation perfored on entering the FilePlanStatus event.
func (p *FilePlanStatusAction) Execute(eventCtx EventContext) EventType {
	resource := eventCtx.(api.Manager)
	resourceSpec := resource.GetSpec()
	fileResourceSpec := resourceSpec.(FileSpec)

	if file.Exists(resource.GetFs(), fileResourceSpec.Path) {
		return FilePlanExists
	}

	return FilePlanDoNotExists
}

// FilePlanNoExistsAction plan the file does not exist.
type FilePlanNoExistsAction struct{}

// Execute implementation perfored on entering the FilePlanNoExists event.
func (p *FilePlanNoExistsAction) Execute(eventCtx EventContext) EventType {
	resource := eventCtx.(api.Manager)
	resource.SetStatus(api.Pending)
	resource.SetMessage("file does not exist")

	return NoOp
}

// FilePlanExistsAction represents the action executed on entering the Exists state.
type FilePlanExistsAction struct{}

// Execute implementation perfored on entering the FilePlanExistss action.
func (p *FilePlanExistsAction) Execute(eventCtx EventContext) EventType {
	resource := eventCtx.(api.Manager)

	resource.SetStatus(api.Pending)
	resource.SetMessage("file exists")

	return NoOp
}
