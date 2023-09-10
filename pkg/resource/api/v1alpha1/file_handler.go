package v1alpha1

import (
	"fmt"

	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
)

// fileHandler handler to manage file changes.
func (f *File) fileHandler() {
	f.doFileMode()

	return
}

// doFileMode implementation to change file mode.
func (f *File) doFileMode() {
	fileMode, err := file.GetMode(f.appFs, f.Spec.GetPath())
	if err != nil {
		f.SetStatusCondition(
			ModeAction, api.Failed, err.Error(), "Unknown", f.Spec.GetModeString())

		return
	}

	fileModeString := fmt.Sprintf("0%o", fileMode.Perm())

	if f.plan {
		// modes differ
		if fileMode != f.Spec.GetMode() {
			f.SetStatusCondition(
				ModeAction,
				api.Pending,
				"modes differ",
				fileModeString,
				f.Spec.GetModeString(),
			)

			return
		}

		// modes are the same
		f.SetStatusCondition(
			ModeAction, api.NoOp, "modes same", fileModeString, f.Spec.GetModeString())

		return
	}

	// modes difer
	if fileMode != f.Spec.GetMode() {
		if err := file.SetMode(f.appFs, f.Spec.GetPath(), f.Spec.GetMode()); err != nil {
			f.SetStatusCondition(
				ModeAction, api.Failed, err.Error(), "Unknown", f.Spec.GetModeString())

			return
		}

		f.SetStatusCondition(
			ModeAction,
			api.Succeeded,
			"modes updated",
			f.Spec.GetModeString(),
			f.Spec.GetModeString(),
		)

		return
	}

	// modes are the same
	f.SetStatusCondition(
		ModeAction, api.NoOp, "modes same", fileModeString, f.Spec.GetModeString())

	return
}
