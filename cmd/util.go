package cmd

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"

	"sigs.k8s.io/yaml"

	"github.com/retr0h/psion/internal/config"
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

	runtimeConfig, err := config.LoadRuntimeConfig(fileContent)
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

	var resourceKind api.Manager = v1alpha1.NewFile(
		logger,
		appFs,
		plan,
	)

	if err := yaml.Unmarshal(fileContent, resourceKind); err != nil {
		return nil, err
	}

	return resourceKind, nil
}

func getAllEmbeddedResourceFiles() ([]*ResourceFilesInfo, error) {
	var files []*ResourceFilesInfo
	if err := fs.WalkDir(eFs, ".", func(filePath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		sHA256String, err := hashEmbededFile(eFs, filePath)
		if err != nil {
			return err
		}

		resourceFilesInfo := &ResourceFilesInfo{
			Path:     filePath,
			Checksum: sHA256String,
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
