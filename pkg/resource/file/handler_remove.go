// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package file

import (
	"github.com/retr0h/psion/pkg/resource/api"
)

// fileRemoveHandler handler to manage file removal.
func (f *File) fileRemoveHandler() {
	if f.file.Exists(f.Spec.GetPath()) {
		f.doFileRemove()

		return
	}
	f.noFileRemove()
}

// doFileRemove implementation to remove file.
func (f *File) doFileRemove() {
	if f.plan {
		f.SetStatusCondition(
			RemoveAction, api.Pending, "file exists", "exists true", "exists false")

		return
	}

	if err := f.file.Remove(f.Spec.GetPath()); err != nil {
		f.SetStatusCondition(
			RemoveAction, api.Failed, err.Error(), "Unknown", "file removed")

		return
	}

	f.SetStatusCondition(
		RemoveAction, api.Succeeded, "file removed", "exists false", "exists false")
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
}
