package v1alpha1

import (
	"errors"
	"log/slog"

	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/retr0h/psion/pkg/resource/api"
)

// ErrNotImplemented is the error returned when feature not implemented.
var ErrNotImplemented = errors.New("not implemented")

const (
	// NoOp represents a no-op event.
	NoOp string = "NoOp"
	// Plan represents the changes to make consistent with the desired state.
	Plan string = "Plan"
	// Apply represents the changes to make the desired state.
	Apply string = "Apply"
)

// File enables declarative updates to File.
type File struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	metav1.TypeMeta   `json:",omitempty,inline"`

	// Spec represents specification of the desired File behavior.
	Spec *FileSpec `json:"spec"`
	// Status contains status of the File.
	Status *api.Status `json:"status"`

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

// NewFile create a new instance of File.
func NewFile(
	logger *slog.Logger,
	appFs afero.Fs,
	plan bool,
) *File {
	return &File{
		Status: &api.Status{
			Conditions: make([]api.StatusConditions, 0, 1),
		},
		Spec:   &FileSpec{},
		plan:   plan,
		logger: logger,
		appFs:  appFs,
	}
}

// GetStatus the status property.
func (f *File) GetStatus() api.Phase { return f.Status.Phase }

// SetStatus set the status property.
func (f *File) SetStatus(status api.Phase) { f.Status.Phase = status }

// GetStatusString the status property as a string.
func (f *File) GetStatusString() string { return string(f.Status.Phase) }

// GetStatusConditions the conditions property.
func (f *File) GetStatusConditions() []api.StatusConditions { return f.Status.Conditions }

// SetStatusCondition set the status condition property.
func (f *File) SetStatusCondition(
	statusType string,
	status api.Phase,
	message string,
	reason string,
	got string,
	want string,
) {
	fileStatusConditions := api.StatusConditions{
		Type:    statusType,
		Status:  status,
		Message: message,
		Reason:  reason,
		Got:     got,
		Want:    want,
	}

	conditions := f.Status.Conditions
	conditions = append(conditions, fileStatusConditions)
	f.Status.Conditions = conditions
}

// GetState provide current state after apply.
func (f *File) GetState() *api.Resource {
	return &api.Resource{
		Name:       f.Name,
		Kind:       f.Kind,
		APIVersion: f.APIVersion,
		Phase:      f.GetStatus(),
		Status: &api.Status{
			Phase:      f.GetStatus(),
			Conditions: f.GetStatusConditions(),
		},
	}
}

// Reconcile make consistent with the desired state.
func (f *File) Reconcile() error {
	if !f.Spec.Exists {
		// file should be deleted
		f.fileRemoveHandler()
	} else {
		// existing file should be modified
		return ErrNotImplemented
	}

	f.logger.Info(
		"completed",
		slog.String("Status", f.GetStatusString()),
		slog.String("Kind", f.Kind),
		slog.String("APIVersion", f.APIVersion),
		slog.Group(FileKind,
			slog.String("Path", f.Spec.Path),
			slog.Bool("Exists", f.Spec.Exists),
		),
		slog.Group("Conditions", f.logStatusConditionGroups()...),
	)

	return nil
}

func (f *File) logStatusConditionGroups() []any {
	var logGroups []any
	for _, condition := range f.Status.Conditions {
		group := slog.Group(condition.GetType(),
			slog.String("Type", condition.GetType()),
			slog.String("Status", condition.GetStatusString()),
			slog.String("Message", condition.GetMessage()),
			slog.String("Reason", condition.GetReason()),
			slog.String("Got", condition.GetGot()),
			slog.String("Want", condition.GetWant()),
		)

		logGroups = append(logGroups, group)
	}

	return logGroups
}
