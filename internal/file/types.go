package file

import (
	"github.com/spf13/afero"
)

// File manages file system crud.
type File struct {
	// appFs FileSystem abstraction.
	appFs afero.Fs
}
