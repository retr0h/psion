package v1alpha1

import (
	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/internal/state"
	"github.com/retr0h/psion/pkg/resource/api"
)

const (
	FileDoStatus    state.StateType = "Unknown"
	FileDoExists    state.StateType = "Exists"
	FileDoNotExists state.StateType = "DoesNotExists"

	FilePlanStatus      state.EventType = "FilePlanStatus"
	FilePlanExists      state.EventType = "FilePlanExists"
	FilePlanDoNotExists state.EventType = "FilePlanDoNotExists"
)

func FilePlanRemoveFSM() *state.StateMachine {
	return &state.StateMachine{
		States: state.States{
			state.Defaults: state.State{
				Events: state.Events{
					FilePlanStatus: FileDoStatus,
				},
			},
			FileDoStatus: state.State{
				Action: &FilePlanStatusAction{},
				Events: state.Events{
					FilePlanExists:      FileDoExists,
					FilePlanDoNotExists: FileDoNotExists,
				},
			},
			FileDoExists: state.State{
				Action: &FilePlanExistsAction{},
				Events: state.Events{
					FilePlanExists: FileDoExists,
				},
			},
			FileDoNotExists: state.State{
				Action: &FilePlanNoExistsAction{},
				Events: state.Events{
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
func (p *FilePlanStatusAction) Execute(
	eventCtx state.EventContext,
) state.EventType {
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
func (p *FilePlanNoExistsAction) Execute(
	eventCtx state.EventContext,
) state.EventType {
	resource := eventCtx.(api.Manager)
	resource.SetStatus(api.Pending)
	resource.SetMessage("file does not exist")

	return state.NoOp
}

// FilePlanExistsAction represents the action executed on entering the Exists state.
type FilePlanExistsAction struct{}

// Execute implementation perfored on entering the FilePlanExistss action.
func (p *FilePlanExistsAction) Execute(
	eventCtx state.EventContext,
) state.EventType {
	resource := eventCtx.(api.Manager)

	resource.SetStatus(api.Pending)
	resource.SetMessage("file exists")

	return state.NoOp
}
