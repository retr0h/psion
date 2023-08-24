package v1alpha1

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/retr0h/psion/internal/file"
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
	logger *logrus.Logger
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
	Message string // "UnexpectedAdmissionError"
	// The reason for the condition's last transition.
	Reason string // "UnexpectedAdmissionError"
}

// NewFile create a new instance of File.
func NewFile(
	logger *logrus.Logger,
	appFs afero.Fs,
	plan bool,
) *File {
	return &File{
		plan:   false,
		logger: logger,
		appFs:  appFs,
	}
}

// Reconcile make consistent with the desired state.
func (f *File) Reconcile() error {
	f.logger.WithFields(logrus.Fields{
		"kind":       f.Kind,
		"apiVersion": f.APIVersion,
	}).Info("reconciling")

	// if exists is false, remove the file and do nothing else
	if !f.Spec.Exists {
		f.Status.Message = fmt.Sprintf("remove: %s", f.Spec.Path)
		if file.Exists(f.appFs, f.Spec.Path) {
			if f.plan {
				f.Status.Phase = api.Pending
				f.Status.Reason = "plan"
			} else {
				err := file.Remove(f.appFs, f.Spec.Path)
				if err != nil {
					f.Status.Reason = err.Error()
					f.Status.Phase = api.Failed
				} else {
					f.Status.Phase = api.Succeeded
				}
			}
		} else {
			f.Status.Phase = api.Failed
			f.Status.Reason = "does not exist"
		}
	}

	f.logger.WithFields(logrus.Fields{
		"Message": f.Status.Message,
		"Reason":  f.Status.Reason,
		"Phase":   f.Status.Phase,
	}).Info("completed")

	return nil
}

// GetStatus the status property.
func (f *File) GetStatus() api.Phase { return f.Status.Phase }

// Get the message property.
func (f *File) GetMessage() string { return f.Status.Message }

// GetReason the reason property.
func (f *File) GetReason() string { return f.Status.Reason }

// // Plan preview the changes to be made.
// func (f *FileResource) Plan(resource map[string]interface{}) error {
// 	return nil
// }

// CopyFile copies the source file to the destination, returning if we changed
// the contents.
// func (f *File) CopyFile(src string, dst string) (bool, error) {
// 	// File doesn't exist - copy it
// 	if !file.Exists(f.appFs, dst) {
// 		err := file.Copy(f.appFs, src, dst)
// 		return true, err
// 	}

// 	// // Are the files identical?
// 	// identical, err := file.Identical(src, dst)
// 	// if err != nil {
// 	// 	return false, err
// 	// }

// 	// // If identical no change
// 	// if identical {
// 	// 	return false, err
// 	// }

// 	// // Since they differ we refresh and that's a change
// 	// err = file.Copy(src, dst)
// 	// return true, err

// 	return true, nil
// }
