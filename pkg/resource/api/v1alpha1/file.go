package v1alpha1

import (
	"errors"
	"log/slog"

	"github.com/spf13/afero"

	"github.com/retr0h/psion/pkg/resource/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ErrNotImplemented is the error returned when feature not implemented.
var ErrNotImplemented = errors.New("not implemented")

// File enables declarative updates to File.
type File struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	metav1.TypeMeta   `json:",omitempty,inline"`

	// Spec represents specification of the desired File behavior.
	Spec FileSpec `json:"spec"`

	// Status contains status of the File.
	Status FileStatus

	// logger logger to be used.
	logger *slog.Logger
	// appFs FileSystem abstraction.
	appFs afero.Fs
	// plan preview the changes to be made.
	plan bool
}

// FileSpec is the specification of the desired behavior of the File.
type FileSpec struct {
	// Exists should file be created or deleted.
	Exists bool `json:"exists"`
	// Path to the file creating or deleting.
	Path string `json:"path,omitempty"`
	// +optional
	// Mode     matcher `json:"mode,omitempty"`
	// Owner    matcher `json:"owner,omitempty"`
	// Group    matcher `json:"group,omitempty"`
	// Contents matcher `json:"contents"`
	// Md5      matcher `json:"md5,omitempty"`
	// Sha256   matcher `json:"sha256,omitempty"`
	// Sha512   matcher `json:"sha512,omitempty"`
	// Skip bool `json:"skip,omitempty"`
}

// FileStatus contains status of the File.
type FileStatus struct {
	// Phase sets `phase` as .status.Phase of the resource.
	Phase api.Phase
	// A human readable message indicating details about the transition.
	Message string
	// The reason for the condition's last transition.
	Reason string
}

// NewFile create a new instance of File.
func NewFile(
	logger *slog.Logger,
	appFs afero.Fs,
	plan bool,
) *File {
	return &File{
		plan:   plan,
		logger: logger,
		appFs:  appFs,
	}
}

// GetStatus the status property.
func (f *File) GetStatus() api.Phase { return f.Status.Phase }

// SetStatus set the status property.
func (f *File) SetStatus(status api.Phase) { f.Status.Phase = status }

// GetStatusAsString the status property cast to a string.
func (f *File) GetStatusAsString() string { return string(f.Status.Phase) }

// GetMessage get the message property.
func (f *File) GetMessage() string { return f.Status.Message }

// SetMessage set the message property.
func (f *File) SetMessage(message string) { f.Status.Message = message }

// GetReason the reason property.
func (f *File) GetReason() string { return f.Status.Reason }

// SetReason set the reason property.
func (f *File) SetReason(reason string) { f.Status.Reason = reason }

// Reconcile make consistent with the desired state.
func (f *File) Reconcile() error {
	// User requeted file be deleted
	if !f.Spec.Exists {
		f.fileRemoveHandler()
	} else {
		return ErrNotImplemented
	}

	f.logger.Info(
		"completed",
		slog.String("Message", f.GetMessage()),
		slog.String("Reason", f.GetReason()),
		slog.String("Phase", f.GetStatusAsString()),
		slog.String("Kind", f.Kind),
		slog.String("APIVersion", f.APIVersion),
		slog.Group(FileKind,
			slog.String("Path", f.Spec.Path),
			slog.Bool("Exists", f.Spec.Exists),
		),
	)

	return nil
}
