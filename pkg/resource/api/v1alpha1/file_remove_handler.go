package v1alpha1

import (
	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
)

// fileRemoveHandler handler to manage file removal.
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
		f.SetStatus(api.Pending)
		f.SetStatusCondition(
			"Remove", api.Pending, "file exists", Plan, "file exists", NoOp)

		return
	}

	f.SetStatus(api.Succeeded)
	if err := file.Remove(f.appFs, f.Spec.Path); err != nil {
		f.SetStatus(api.Failed)
		f.SetStatusCondition(
			"Remove", api.Failed, err.Error(), Apply, "file exists", "file removed")
	}

	f.SetStatusCondition(
		"Remove", api.Succeeded, "file removed", Apply, "file exists", "file removed")

	return
}

// noFileRemove implementation to not remove removal.
func (f *File) noFileRemove() {
	if f.plan {
		f.SetStatus(api.Pending)
		f.SetStatusCondition(
			"Remove", api.Pending, "file does not exist", Plan, "file does not exist", NoOp)

		return
	}

	f.SetStatus(api.Succeeded)
	f.SetStatusCondition(
		"Remove", api.Succeeded, "file does not exist", Apply, "file does not exist", NoOp)

	return
}
