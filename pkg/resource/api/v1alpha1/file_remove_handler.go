package v1alpha1

import (
	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
)

// fileRemoveHandler handler to manage file removal.
func (f *File) fileRemoveHandler() {
	if file.Exists(f.appFs, f.Spec.GetPath()) {
		f.doFileRemove()

		return
	}
	f.noFileRemove()

	return
}

// doFileRemove implementation to remove file.
func (f *File) doFileRemove() {
	if f.plan {
		f.SetStatusCondition(
			RemoveAction, api.Pending, "file exists", "exists true", "exists false")

		return
	}

	if err := file.Remove(f.appFs, f.Spec.GetPath()); err != nil {
		f.SetStatusCondition(
			RemoveAction, api.Failed, err.Error(), "Unknown", "file removed")

		return
	}

	f.SetStatusCondition(
		RemoveAction, api.Succeeded, "file removed", "exists false", "exists false")

	return
}

// noFileRemove implementation to not remove removal.
func (f *File) noFileRemove() {
	if f.plan {
		f.SetStatusCondition(
			RemoveAction, api.NoOp, "file does not exist", "exists false", "exists false")

		return
	}

	f.SetStatusCondition(
		RemoveAction, api.NoOp, "file does not exist", "exists false", "exists false")

	return
}
