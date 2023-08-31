// Package file contains some simple utility functions.
package file

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"

	"github.com/spf13/afero"
)

// Read reads the contents of the filePath.
func Read(
	fs afero.Fs,
	filePath string,
) ([]byte, error) {
	f, err := fs.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	fileinfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	filesize := fileinfo.Size()
	buf := make([]byte, filesize)

	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		if n > 0 {
			return buf, nil
		}

	}
}

// Copy copies the contents of the src file to the dst file.
func Copy(
	fs afero.Fs,
	src string,
	dst string,
) error {
	r, err := fs.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()

	w, err := fs.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	return w.Close()
}

// test this TODO
// Remove removes the named file if exists.
func Remove(
	fs afero.Fs,
	filePath string,
) error {
	if err := fs.Remove(filePath); err != nil {
		return err
	}

	return nil
}

// Exists reports if the named file or directory exists.
func Exists(
	fs afero.Fs,
	filePath string,
) bool {
	if _, err := fs.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Size returns the named files size in bytes.
func Size(
	fs afero.Fs,
	filePath string,
) (int64, error) {
	fi, err := fs.Stat(filePath)
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

// HashFile returns the SHA1-hash of the contents of the specified file.
func HashFile(
	fs afero.Fs,
	filePath string,
) (string, error) {
	var returnSHA1String string

	file, err := fs.Open(filePath)
	if err != nil {
		return returnSHA1String, err
	}

	defer func() { _ = file.Close() }()

	hash := sha1.New()

	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA1String, err
	}

	hashInBytes := hash.Sum(nil)[:20]
	returnSHA1String = hex.EncodeToString(hashInBytes)

	return returnSHA1String, nil
}

// Identical compares the contents of the two specified files, returning
// true if they're identical.
func Identical(
	fs afero.Fs,
	a string,
	b string,
) (bool, error) {
	hashA, errA := HashFile(fs, a)
	if errA != nil {
		return false, errA
	}

	hashB, errB := HashFile(fs, b)
	if errB != nil {
		return false, errB
	}

	// Are the hashes are identical?
	// If so then the files are identical.
	if hashA == hashB {
		return true, nil
	}

	return false, nil
}
