package v1alpha1

import (
	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/internal/state"
	"github.com/retr0h/psion/pkg/resource/api"
)

const (
	// FileDoStatus state of file unknown.
	FileDoStatus state.StateType = "Unknown"
	// FileDoExists state of file existing.
	FileDoExists state.StateType = "Exists"
	// FileDoNotExists state of file does not exist.
	FileDoNotExists state.StateType = "DoesNotExists"

	// FilePlanStatusEvent planned file unknown event.
	FilePlanStatusEvent state.EventType = "FilePlanStatus"
	// FilePlanExistsEvent planned file existing event.
	FilePlanExistsEvent state.EventType = "FilePlanExists"
	// FilePlanDoNotExistsEvent planned file does not exist event.
	FilePlanDoNotExistsEvent state.EventType = "FilePlanDoNotExists"
	// FileStatusEvent file unknown event.
	FileStatusEvent state.EventType = "FileStatus"
	// FileExistsEvent file existing event.
	FileExistsEvent state.EventType = "FileExists"
	// FileDoNotExistsEvent file does not exist event.
	FileDoNotExistsEvent state.EventType = "FileDoNotExists"
)

// FilePlanRemoveFSM finite state machine planning file remove.
func FilePlanRemoveFSM() *state.StateMachine {
	return &state.StateMachine{
		States: state.States{
			state.Defaults: state.State{
				Events: state.Events{
					FilePlanStatusEvent: FileDoStatus,
				},
			},
			FileDoStatus: state.State{
				Action: &FilePlanStatusAction{},
				Events: state.Events{
					FilePlanExistsEvent:      FileDoExists,
					FilePlanDoNotExistsEvent: FileDoNotExists,
				},
			},
			FileDoExists: state.State{
				Action: &FilePlanExistsAction{},
				Events: state.Events{
					FilePlanExistsEvent: FileDoExists,
				},
			},
			FileDoNotExists: state.State{
				Action: &FilePlanNoExistsAction{},
				Events: state.Events{
					FilePlanDoNotExistsEvent: FileDoNotExists,
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
		return FilePlanExistsEvent
	}

	return FilePlanDoNotExistsEvent
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

// FileRemoveFSM finite state machine file remove.
func FileRemoveFSM() *state.StateMachine {
	return &state.StateMachine{
		States: state.States{
			state.Defaults: state.State{
				Events: state.Events{
					FileStatusEvent: FileDoStatus,
				},
			},
			FileDoStatus: state.State{
				Action: &FileStatusAction{},
				Events: state.Events{
					FileExistsEvent:      FileDoExists,
					FileDoNotExistsEvent: FileDoNotExists,
				},
			},
			FileDoExists: state.State{
				Action: &FileExistsAction{},
				Events: state.Events{
					FileExistsEvent: FileDoExists,
				},
			},
			FileDoNotExists: state.State{
				Action: &FileNoExistsAction{},
				Events: state.Events{
					FileDoNotExistsEvent: FileDoNotExists,
				},
			},
		},
	}
}

// FileStatusAction initial state when entering the state machine when
// Plan is not set and the File resource should be removed.
type FileStatusAction struct{}

// Execute implementation perfored on entering the FileStatus event.
func (p *FileStatusAction) Execute(
	eventCtx state.EventContext,
) state.EventType {
	resource := eventCtx.(api.Manager)
	resourceSpec := resource.GetSpec()
	fileResourceSpec := resourceSpec.(FileSpec)

	if file.Exists(resource.GetFs(), fileResourceSpec.Path) {
		return FileExistsEvent
	}

	return FileDoNotExistsEvent
}

// FileNoExistsAction the file does not exist.
type FileNoExistsAction struct{}

// Execute implementation perfored on entering the FileNoExists event.
func (p *FileNoExistsAction) Execute(
	eventCtx state.EventContext,
) state.EventType {
	resource := eventCtx.(api.Manager)
	resource.SetStatus(api.Succeeded)
	resource.SetMessage("file does not exist")

	return state.NoOp
}

// FileExistsAction represents the action executed on entering the Exists state.
type FileExistsAction struct{}

// Execute implementation perfored on entering the FileExists action.
func (p *FileExistsAction) Execute(
	eventCtx state.EventContext,
) state.EventType {
	resource := eventCtx.(api.Manager)
	resourceSpec := resource.GetSpec()
	fileResourceSpec := resourceSpec.(FileSpec)

	resource.SetStatus(api.Succeeded)
	resource.SetMessage("file removed")

	if err := file.Remove(resource.GetFs(), fileResourceSpec.Path); err != nil {
		resource.SetStatus(api.Failed)
		resource.SetMessage(err.Error())

		return state.NoOp
	}

	return state.NoOp
}
