package file

import (
	"log/slog"

	"github.com/retr0h/psion/pkg/resource/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/status"
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
	Status *status.Status `json:"status"`

	// logger logger to be used.
	logger *slog.Logger
	// appFs FileSystem abstraction.
	// appFs afero.Fs
	// plan preview the changes to be made.
	plan bool
	// file manager repository.
	file internal.FileManager
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
