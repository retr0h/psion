package cmd

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"

	"github.com/retr0h/psion/internal/config"
	"github.com/retr0h/psion/pkg/resource/api"
	"github.com/retr0h/psion/pkg/resource/api/v1alpha1"
)

func loadResourceFile(
	fs afero.Fs,
	efs embed.FS,
	filePath string,
	plan bool,
) (api.Manager, error) {
	// read from the embedded fs
	fileContent, err := efs.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	runtimeConfig, err := config.LoadRuntimeConfig(fileContent)
	if err != nil {
		return nil, fmt.Errorf("cannot load file: %w", err)
	}

	if runtimeConfig.APIVersion != v1alpha1.FileAPIVersion {
		return nil, fmt.Errorf("invalid apiVersion: %s file: %s", runtimeConfig.APIVersion, filePath)
	}

	// currently only support the File Kind
	if runtimeConfig.Kind != v1alpha1.FileKind {
		return nil, fmt.Errorf("invalid kind: %s file: %s", runtimeConfig.Kind, filePath)
	}

	var resourceKind api.Manager = v1alpha1.NewFile(
		logger,
		fs,
		plan,
	)

	if err := yaml.Unmarshal(fileContent, resourceKind); err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %w", err)
	}

	return resourceKind, nil
}

func getAllEmbeddedResourceFiles(efs embed.FS) ([]string, error) {
	var files []string
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func loadAllEmbeddedResourceFiles(
	fs afero.Fs,
	efs embed.FS,
	plan bool,
) ([]api.Manager, error) {
	files, err := getAllEmbeddedResourceFiles(resourceFiles)
	if err != nil {
		return nil, fmt.Errorf("cannot walk dir: %w", err)
	}

	resources := make([]api.Manager, 0, 1)
	for _, filePath := range files {
		resourceFile, err := loadResourceFile(fs, efs, filePath, plan)
		if err != nil {
			return nil, fmt.Errorf("cannot load resource file: %w", err)
		}
		resources = append(resources, resourceFile)
	}

	return resources, nil
}
