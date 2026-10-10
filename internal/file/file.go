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

// Package file contains some simple utility functions.
package file

import (
	"io"
	"io/fs"
	"os"

	"github.com/spf13/afero"
)

// New create a file instance.
func New(
	appFs afero.Fs,
) *File {
	return &File{
		appFs: appFs,
	}
}

// Read reads the contents of the filePath.
func (f *File) Read(
	filePath string,
) ([]byte, error) {
	file, err := f.appFs.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	filesize := fileInfo.Size()
	buf := make([]byte, filesize)

	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		if n > 0 {
			return buf, nil
		}

	}
}

// Remove removes the named file if exists.
func (f *File) Remove(
	filePath string,
) error {
	return f.appFs.Remove(filePath)
}

// Exists reports if the named file or directory exists.
func (f *File) Exists(
	filePath string,
) bool {
	if _, err := f.appFs.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// GetMode returns the named files mode.
func (f *File) GetMode(
	filePath string,
) (fs.FileMode, error) {
	fileInfo, err := f.appFs.Stat(filePath)
	if err != nil {
		return 0, err
	}

	return fileInfo.Mode(), nil
}

// SetMode sets the named files mode.
func (f *File) SetMode(
	filePath string,
	mode fs.FileMode,
) error {
	return f.appFs.Chmod(filePath, mode)
}

// // Copy copies the contents of the src file to the dst file.
