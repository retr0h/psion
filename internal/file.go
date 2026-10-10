package internal

import (
	"io/fs"
)

// FileManager manager responsible for File operations.
type FileManager interface {
	Read(
		filePath string,
	) ([]byte, error)
	Remove(
		filePath string,
	) error
	Exists(
		filePath string,
	) bool
	GetMode(
		filePath string,
	) (fs.FileMode, error)
	SetMode(
		filePath string,
		mode fs.FileMode,
	) error
}
