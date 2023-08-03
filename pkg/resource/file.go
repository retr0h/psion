package resource

import (
	"github.com/spf13/afero"

	"github.com/retr0h/psion/internal/file"
)

type FileManager interface {
	GetAuthToken() (string, error)
	GetConfiguration(string) (string, error)
}

type File struct {
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// Meta     meta    `json:"meta,omitempty" yaml:"meta,omitempty"`
	id   string `json:"-" yaml:"-"`
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	// Exists   matcher `json:"exists" yaml:"exists"`
	// Mode     matcher `json:"mode,omitempty" yaml:"mode,omitempty"`
	// Size     matcher `json:"size,omitempty" yaml:"size,omitempty"`
	// Owner    matcher `json:"owner,omitempty" yaml:"owner,omitempty"`
	// Group    matcher `json:"group,omitempty" yaml:"group,omitempty"`
	// LinkedTo matcher `json:"linked-to,omitempty" yaml:"linked-to,omitempty"`
	// Filetype matcher `json:"filetype,omitempty" yaml:"filetype,omitempty"`
	// Contains matcher `json:"contains,omitempty" yaml:"contains,omitempty"`
	// Contents matcher `json:"contents" yaml:"contents"`
	// Md5      matcher `json:"md5,omitempty" yaml:"md5,omitempty"`
	// Sha256   matcher `json:"sha256,omitempty" yaml:"sha256,omitempty"`
	// Sha512   matcher `json:"sha512,omitempty" yaml:"sha512,omitempty"`
	Skip bool `json:"skip,omitempty" yaml:"skip,omitempty"`

	appFs afero.Fs
	//= afero.NewOsFs()
}

func NewFile() (*File, error) {
	return nil, nil
}

// func (f *File) GetTitle() string { return f.Title }
// func (f *File) GetMeta() meta    { return f.Meta }

// CopyFile copies the source file to the destination, returning if we changed
// the contents.
func (f *File) CopyFile(src string, dst string) (bool, error) {
	// File doesn't exist - copy it
	if !file.Exists(f.appFs, dst) {
		err := file.Copy(f.appFs, src, dst)
		return true, err
	}

	// // Are the files identical?
	// identical, err := file.Identical(src, dst)
	// if err != nil {
	// 	return false, err
	// }

	// // If identical no change
	// if identical {
	// 	return false, err
	// }

	// // Since they differ we refresh and that's a change
	// err = file.Copy(src, dst)
	// return true, err

	return true, nil
}
