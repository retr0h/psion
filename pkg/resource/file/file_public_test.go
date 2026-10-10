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

package file_test

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/yaml"

	intfile "github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
	resource "github.com/retr0h/psion/pkg/resource/file"
)

type FilePublicTestSuite struct {
	suite.Suite
}

func (s *FilePublicTestSuite) TestReconcile() {
	tests := []struct {
		name          string
		plan          bool
		desiredExists bool
		exists        bool
		mode          fs.FileMode
		readOnly      bool
		wantPhase     api.Phase
		wantType      api.SpecAction
		wantMessage   string
		wantGot       string
		wantWant      string
		wantExists    bool
		wantMode      fs.FileMode
	}{
		{
			name: "plan mode change", plan: true, desiredExists: true,
			exists: true, mode: 0o644,
			wantPhase: api.Pending, wantType: resource.ModeAction,
			wantMessage: "modes differ", wantGot: "0644", wantWant: "0700",
			wantExists: true, wantMode: 0o644,
		},
		{
			name: "plan matching mode", plan: true, desiredExists: true,
			exists: true, mode: 0o700,
			wantPhase: api.NoOp, wantType: resource.ModeAction,
			wantMessage: "modes same", wantGot: "0700", wantWant: "0700",
			wantExists: true, wantMode: 0o700,
		},
		{
			name: "plan missing file mode", plan: true, desiredExists: true,
			wantPhase: api.Failed, wantType: resource.ModeAction,
			wantMessage: "open /app/filePath: file does not exist", wantGot: "Unknown", wantWant: "0700",
		},
		{
			name: "apply mode change", desiredExists: true,
			exists: true, mode: 0o644,
			wantPhase: api.Succeeded, wantType: resource.ModeAction,
			wantMessage: "modes updated", wantGot: "0700", wantWant: "0700",
			wantExists: true, wantMode: 0o700,
		},
		{
			name: "apply matching mode", desiredExists: true,
			exists: true, mode: 0o700,
			wantPhase: api.NoOp, wantType: resource.ModeAction,
			wantMessage: "modes same", wantGot: "0700", wantWant: "0700",
			wantExists: true, wantMode: 0o700,
		},
		{
			name: "apply missing file mode", desiredExists: true,
			wantPhase: api.Failed, wantType: resource.ModeAction,
			wantMessage: "open /app/filePath: file does not exist", wantGot: "Unknown", wantWant: "0700",
		},
		{
			name: "apply mode change denied", desiredExists: true,
			exists: true, mode: 0o644, readOnly: true,
			wantPhase: api.Failed, wantType: resource.ModeAction,
			wantMessage: "operation not permitted", wantGot: "Unknown", wantWant: "0700",
			wantExists: true, wantMode: 0o644,
		},
		{
			name: "plan removal", plan: true, exists: true, mode: 0o644,
			wantPhase: api.Pending, wantType: resource.RemoveAction,
			wantMessage: "file exists", wantGot: "exists true", wantWant: "exists false",
			wantExists: true, wantMode: 0o644,
		},
		{
			name: "plan already absent", plan: true,
			wantPhase: api.NoOp, wantType: resource.RemoveAction,
			wantMessage: "file does not exist", wantGot: "exists false", wantWant: "exists false",
		},
		{
			name: "apply removal", exists: true, mode: 0o644,
			wantPhase: api.Succeeded, wantType: resource.RemoveAction,
			wantMessage: "file removed", wantGot: "exists false", wantWant: "exists false",
		},
		{
			name:      "apply already absent",
			wantPhase: api.NoOp, wantType: resource.RemoveAction,
			wantMessage: "file does not exist", wantGot: "exists false", wantWant: "exists false",
		},
		{
			name: "apply removal denied", exists: true, mode: 0o644, readOnly: true,
			wantPhase: api.Failed, wantType: resource.RemoveAction,
			wantMessage: "operation not permitted", wantGot: "Unknown", wantWant: "file removed",
			wantExists: true, wantMode: 0o644,
		},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			const filePath = "/app/filePath"
			appFs := afero.NewMemMapFs()
			s.Require().NoError(appFs.MkdirAll("/app", 0o755))
			if tc.exists {
				s.Require().
					NoError(afero.WriteFile(appFs, filePath, []byte("mockContent"), tc.mode))
			}
			if tc.readOnly {
				appFs = afero.NewReadOnlyFs(appFs)
			}
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			instance := resource.NewFile(logger, intfile.New(appFs), tc.plan)
			content := fmt.Sprintf(`
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: name
spec:
  exists: %t
  path: %s
`, tc.desiredExists, filePath)
			if tc.desiredExists {
				content += "  mode: 0o700\n"
			}
			s.Require().NoError(yaml.Unmarshal([]byte(content), instance))

			s.Require().NoError(instance.Reconcile())
			s.Equal(tc.wantPhase, instance.GetStatus())
			conditions := instance.GetStatusConditions()
			s.Require().Len(conditions, 1)
			reason := api.Apply
			if tc.plan {
				reason = api.Plan
			}
			s.Equal(api.StatusConditions{
				Type:    tc.wantType,
				Status:  tc.wantPhase,
				Message: tc.wantMessage,
				Reason:  reason,
				Got:     tc.wantGot,
				Want:    tc.wantWant,
			}, conditions[0])

			exists, err := afero.Exists(appFs, filePath)
			s.Require().NoError(err)
			s.Equal(tc.wantExists, exists)
			if tc.wantExists {
				info, err := appFs.Stat(filePath)
				s.Require().NoError(err)
				s.Equal(tc.wantMode, info.Mode())
			}
		})
	}
}

func TestFilePublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(FilePublicTestSuite))
}
