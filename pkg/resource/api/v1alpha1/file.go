package v1alpha1

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	// "github.com/retr0h/psion/internal/file"
)

// File enables declarative updates to Files.
type File struct {
	// Standard object metadata.
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	// Spec represents specification of the desired File behavior
	Spec FileSpec `json:"spec" yaml:"spec"`

	// apiVersion version of API resource to use.
	apiVersion string
	// kind is the type of resource being referenced.
	kind string
	// logger
	logger *logrus.Logger
	// appFs FileSystem abstraction.
	appFs afero.Fs
}

// ObjectMeta is metadata that all resources must have.
type ObjectMeta struct {
	// Name must be unique within a resource.
	Name string `json:"name" yaml:"name"`

	// Annotations is an unstructured key value map stored with a resource to
	// store and retrieve arbitrary metadata.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations"`
}

// FileSpec is the specification of the desired behavior of the File.
type FileSpec struct {
	// Exists should file be created or deleted.
	Exists bool `json:"exists" yaml:"exists"`
	// Path   string `json:"path,omitempty" yaml:"path,omitempty"`
	// +optional
	// Mode     matcher `json:"mode,omitempty" yaml:"mode,omitempty"`
	// Owner    matcher `json:"owner,omitempty" yaml:"owner,omitempty"`
	// Group    matcher `json:"group,omitempty" yaml:"group,omitempty"`
	// Contents matcher `json:"contents" yaml:"contents"`
	// Md5      matcher `json:"md5,omitempty" yaml:"md5,omitempty"`
	// Sha256   matcher `json:"sha256,omitempty" yaml:"sha256,omitempty"`
	// Sha512   matcher `json:"sha512,omitempty" yaml:"sha512,omitempty"`
	// Skip bool `json:"skip,omitempty" yaml:"skip,omitempty"`
}

// NewFile create a new instance of File.
func NewFile(
	apiVersion string,
	kind string,
	logger *logrus.Logger,
	appFs afero.Fs,
	resource []byte,
) *File {
	return &File{
		apiVersion: apiVersion,
		kind:       kind,
		logger:     logger,
		appFs:      appFs,
	}
}

// Reconcile make consistent with the desired state.
func (f *File) Reconcile() error {
	f.logger.WithFields(logrus.Fields{
		"kind":       f.kind,
		"apiVersion": f.apiVersion,
	}).Info("reconciling")

	return nil
}

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
