package file_test

import (
	"github.com/spf13/afero"
)

type FileSpec struct {
	appFs    afero.Fs
	srcDir   string
	srcFile  string
	srcFiles []string
}

func createFileSpecs(specs []FileSpec) {
	for _, s := range specs {
		if s.srcDir != "" {
			_ = s.appFs.MkdirAll(s.srcDir, 0o755)
		}

		if s.srcFile != "" {
			_ = afero.WriteFile(
				s.appFs,
				s.srcFile,
				[]byte("mockContent"),
				0o644,
			)
		}

		if len(s.srcFiles) > 0 {
			for _, f := range s.srcFiles {
				_, _ = s.appFs.Create(f)
			}
		}
	}
}
