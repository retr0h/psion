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
	"log/slog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
)

const (
	// Kind resource used to manage this resource.
	Kind = "File"
	// APIVersion current version of File API.
	APIVersion = "files.psion.io/v1alpha1"

	// RemoveAction action type.
	RemoveAction api.SpecAction = "Remove"
	// ModeAction action type.
	ModeAction api.SpecAction = "Mode"
)

// File enables declarative updates to File.
type File struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	metav1.TypeMeta   `json:",omitempty,inline"`

	// Spec represents specification of the desired File behavior.
	Spec *Spec `json:"spec"`
	// Status contains status of the File.
	Status *api.Status `json:"status"`

	// logger logger to be used.
	logger *slog.Logger
	// appFs FileSystem abstraction.
	// appFs afero.Fs
	// plan preview the changes to be made.
	plan bool
	// file manager repository.
	file *file.File
}

// Spec is the specification of the desired behavior of the File.
type Spec struct {
	// Exists should file be created or deleted.
	Exists bool `json:"exists"`
	// Path to the file creating or deleting.
	Path string `json:"path,omitempty"`
	// +optional
	// Mode the mode (permissions) the file should be.
	Mode int `json:"mode,omitempty"`
	// Owner    matcher `json:"owner,omitempty"`
	// Group    matcher `json:"group,omitempty"`
	// Contents matcher `json:"contents"`
	// Md5      matcher `json:"md5,omitempty"`
	// Sha256   matcher `json:"sha256,omitempty"`
	// Sha512   matcher `json:"sha512,omitempty"`
	// Skip bool `json:"skip,omitempty"`
}
