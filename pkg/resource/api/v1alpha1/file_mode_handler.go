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

package v1alpha1

import (
	"fmt"

	"github.com/retr0h/psion/pkg/resource/api"
)

// doFileMode implementation to change file mode.
func (f *File) doFileMode() {
	fileMode, err := f.file.GetMode(f.Spec.GetPath())
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
		if err := f.file.SetMode(f.Spec.GetPath(), f.Spec.GetMode()); err != nil {
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
}
