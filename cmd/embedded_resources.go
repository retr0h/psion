// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package cmd

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"

	"sigs.k8s.io/yaml"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/config"
	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
	"github.com/retr0h/psion/pkg/resource/api/v1alpha1"
)

func loadResourceFile(
	filePath string,
	plan bool,
) (api.Manager, error) {
	// read from the embedded fs
	fileContent, err := eFs.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var c internal.ConfigManager = config.New()
	runtimeConfig, err := c.GetConfig(fileContent)
	if err != nil {
		return nil, err
	}

	if runtimeConfig.APIVersion != v1alpha1.FileAPIVersion {
		return nil, fmt.Errorf(
			"invalid apiVersion: %s file: %s",
			runtimeConfig.APIVersion,
			filePath,
		)
	}

	// currently only support the File Kind
	if runtimeConfig.Kind != v1alpha1.FileKind {
		return nil, fmt.Errorf("invalid kind: %s file: %s", runtimeConfig.Kind, filePath)
	}

	fileManager := file.New(appFs)
	var resourceKind api.Manager = v1alpha1.NewFile(
		logger,
		fileManager,
		plan,
	)

	if err := yaml.Unmarshal(fileContent, resourceKind); err != nil {
		return nil, err
	}

	return resourceKind, nil
}

func getAllEmbeddedResourceFiles() ([]*ResourceFilesInfo, error) {
	var files []*ResourceFilesInfo
	if err := fs.WalkDir(eFs, ".", func(filePath string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			return nil
		}

		checksum := "unknown"
		if sHA256String, _ := hashEmbededFile(eFs, filePath); sHA256String != "" {
			checksum = sHA256String
		}

		resourceFilesInfo := &ResourceFilesInfo{
			Path:     filePath,
			Checksum: checksum,
			Type:     "SHA256",
		}

		files = append(files, resourceFilesInfo)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func loadAllEmbeddedResourceFiles(
	plan bool,
) ([]api.Manager, error) {
	files, err := getAllEmbeddedResourceFiles()
	if err != nil {
		return nil, err
	}

	resources := make([]api.Manager, 0, 1)
	for _, resourceFileInfo := range files {
		resourceFile, err := loadResourceFile(resourceFileInfo.Path, plan)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resourceFile)
	}

	return resources, nil
}

func hashEmbededFile(
	eFs embed.FS,
	filePath string,
) (string, error) {
	var returnSHA256String string

	file, err := eFs.Open(filePath)
	if err != nil {
		return returnSHA256String, err
	}

	defer func() { _ = file.Close() }()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	returnSHA256String = hex.EncodeToString(hash.Sum(nil))

	return returnSHA256String, nil
}
