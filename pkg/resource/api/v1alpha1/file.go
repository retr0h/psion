package v1alpha1

import (
	"fmt"
	"log/slog"

	"github.com/spf13/afero"

	"github.com/retr0h/psion/pkg/resource/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// File enables declarative updates to File.
type File struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	metav1.TypeMeta   `json:",omitempty,inline"`

	// Spec represents specification of the desired File behavior
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

// Set the status property.
func (f *File) SetStatus(status api.Phase) { f.Status.Phase = status }

// GetStatusAsString the status property cast to a string.
func (f *File) GetStatusAsString() string { return string(f.Status.Phase) }

// Get the message property.
func (f *File) GetMessage() string { return f.Status.Message }

// Set the message property.
func (f *File) SetMessage(message string) { f.Status.Message = message }

// GetReason the reason property.
func (f *File) GetReason() string { return f.Status.Reason }

// GetSpec the spec property.
func (f *File) GetSpec() interface{} { return f.Spec }

// GetFs the appFs property.
func (f *File) GetFs() afero.Fs { return f.appFs }

// Reconcile make consistent with the desired state.
func (f *File) Reconcile() error {
	f.logger.Info(
		"reconciling",
		slog.String("Kind", f.Kind),
		slog.String("APIVersion", f.APIVersion),
	)

	// Spec.File.Path should be removed
	if !f.Spec.Exists {
		if f.plan {
			// Plan the removal
			f.Status.Reason = "Plan"
			planFSM := FilePlanRemoveFSM()
			err := planFSM.SendEvent(FilePlanStatus, f)
			if err != nil {
				fmt.Errorf("Couldn't set the initial state of the state machine, err: %w", err)
			}
		} else {
			// Do the removal
			fmt.Println("IMPLEMENT")
		}
	}

	f.logger.Info(
		"completed",
		slog.String("Message", f.GetMessage()),
		slog.String("Reason", f.GetReason()),
		slog.String("Phase", f.GetStatusAsString()),
		slog.Group(FileKind,
			slog.String("Path", f.Spec.Path),
			slog.Bool("Exists", f.Spec.Exists),
		),
	)

	return nil
}
