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
	"io/fs"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/psion/internal/file"
)

type FilePublicTestSuite struct {
	suite.Suite
}

func (s *FilePublicTestSuite) TestRead() {
	tests := []struct {
		name    string
		exists  bool
		content string
		wantErr error
	}{
		{name: "existing file", exists: true, content: "mockContent"},
		{name: "missing file", wantErr: fs.ErrNotExist},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			appFs := afero.NewMemMapFs()
			if tc.exists {
				s.Require().NoError(afero.WriteFile(appFs, "/file.txt", []byte(tc.content), 0o644))
			}
			got, err := file.New(appFs).Read("/file.txt")
			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
				s.Nil(got)
				return
			}
			s.Require().NoError(err)
			s.Equal(tc.content, string(got))
		})
	}
}

func (s *FilePublicTestSuite) TestRemove() {
	tests := []struct {
		name       string
		exists     bool
		readOnly   bool
		wantErr    error
		wantExists bool
	}{
		{name: "existing file", exists: true},
		{name: "missing file", wantErr: fs.ErrNotExist},
		{
			name:       "read only filesystem",
			exists:     true,
			readOnly:   true,
			wantErr:    fs.ErrPermission,
			wantExists: true,
		},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			appFs := afero.NewMemMapFs()
			if tc.exists {
				s.Require().
					NoError(afero.WriteFile(appFs, "/file.txt", []byte("mockContent"), 0o644))
			}
			if tc.readOnly {
				appFs = afero.NewReadOnlyFs(appFs)
			}
			err := file.New(appFs).Remove("/file.txt")
			s.ErrorIs(err, tc.wantErr)
			exists, err := afero.Exists(appFs, "/file.txt")
			s.Require().NoError(err)
			s.Equal(tc.wantExists, exists)
		})
	}
}

func (s *FilePublicTestSuite) TestExists() {
	tests := []struct {
		name   string
		exists bool
	}{
		{name: "existing file", exists: true},
		{name: "missing file"},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			appFs := afero.NewMemMapFs()
			if tc.exists {
				s.Require().
					NoError(afero.WriteFile(appFs, "/file.txt", []byte("mockContent"), 0o644))
			}
			s.Equal(tc.exists, file.New(appFs).Exists("/file.txt"))
		})
	}
}

func (s *FilePublicTestSuite) TestGetMode() {
	tests := []struct {
		name    string
		exists  bool
		mode    fs.FileMode
		wantErr error
	}{
		{name: "regular file", exists: true, mode: 0o644},
		{name: "executable file", exists: true, mode: 0o755},
		{name: "missing file", wantErr: fs.ErrNotExist},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			appFs := afero.NewMemMapFs()
			if tc.exists {
				s.Require().
					NoError(afero.WriteFile(appFs, "/file.txt", []byte("mockContent"), tc.mode))
			}
			got, err := file.New(appFs).GetMode("/file.txt")
			s.ErrorIs(err, tc.wantErr)
			s.Equal(tc.mode, got)
		})
	}
}

func (s *FilePublicTestSuite) TestSetMode() {
	tests := []struct {
		name     string
		exists   bool
		readOnly bool
		wantErr  error
		wantMode fs.FileMode
	}{
		{name: "existing file", exists: true, wantMode: 0o700},
		{name: "missing file", wantErr: fs.ErrNotExist},
		{
			name:     "read only filesystem",
			exists:   true,
			readOnly: true,
			wantErr:  fs.ErrPermission,
			wantMode: 0o644,
		},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			appFs := afero.NewMemMapFs()
			if tc.exists {
				s.Require().
					NoError(afero.WriteFile(appFs, "/file.txt", []byte("mockContent"), 0o644))
			}
			if tc.readOnly {
				appFs = afero.NewReadOnlyFs(appFs)
			}
			err := file.New(appFs).SetMode("/file.txt", 0o700)
			s.ErrorIs(err, tc.wantErr)
			if tc.exists {
				info, err := appFs.Stat("/file.txt")
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
