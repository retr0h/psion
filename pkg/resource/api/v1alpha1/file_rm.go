package v1alpha1

import (
	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
)

// doFileRemoveHandler handler to manage file removal.
func (f *File) fileRemoveHandler() {
	if file.Exists(f.appFs, f.Spec.Path) {
		f.doFileRemove()

		return
	}
	f.noFileRemove()

	return
}

// doFileRemove implementation to remove file.
func (f *File) doFileRemove() {
	if f.plan {
		f.SetReason("Plan")
		f.SetStatus(api.Pending)
		f.SetMessage("file exists")

		return
	}

	f.SetReason("Apply")
	f.SetStatus(api.Succeeded)
	f.SetMessage("file removed")

	if err := file.Remove(f.appFs, f.Spec.Path); err != nil {
		f.SetStatus(api.Failed)
		f.SetMessage(err.Error())
	}

	return
}

// noFileRemove implementation to not remove removal.
func (f *File) noFileRemove() {
	if f.plan {
		f.SetReason("Plan")
		f.SetStatus(api.Pending)
		f.SetMessage("file does not exist")

		return
	}

	f.SetReason("Apply")
	f.SetStatus(api.Succeeded)
	f.SetMessage("file does not exist")

	return
}
