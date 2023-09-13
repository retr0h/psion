package file

import (
	"github.com/spf13/afero"
)

// File manages file crud.
type File struct {
	// appFs FileSystem abstraction.
	appFs afero.Fs
}
