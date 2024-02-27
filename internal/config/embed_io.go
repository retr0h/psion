package config

import (
	"fmt"
	// "crypto/sha256"
	// "embed"
	// "encoding/hex"
	// "io"
	"io/fs"
	// "sigs.k8s.io/yaml"
	// "github.com/retr0h/psion/internal"
	// "github.com/retr0h/psion/internal/config"
	// intFile "github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/api"
	// "github.com/retr0h/psion/pkg/resource/file"
)

// func loadResourceFile(
// 	filePath string,
// 	plan bool,
// ) (resource.Manager, error) {
// 	// read from the embedded fs
// 	fileContent, err := eFs.ReadFile(filePath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var c internal.ConfigManager = config.New()
// 	runtimeConfig, err := c.GetConfig(fileContent)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if runtimeConfig.APIVersion != file.APIVersion {
// 		return nil, fmt.Errorf(
// 			"invalid apiVersion: %s file: %s",
// 			runtimeConfig.APIVersion,
// 			filePath,
// 		)
// 	}

// 	// currently only support the File Kind
// 	if runtimeConfig.Kind != file.Kind {
// 		return nil, fmt.Errorf("invalid kind: %s file: %s", runtimeConfig.Kind, filePath)
// 	}

// 	fileManager := intFile.New(appFs)
// 	var resourceKind resource.Manager = file.New(
// 		logger,
// 		fileManager,
// 		plan,
// 	)

// 	if err := yaml.Unmarshal(fileContent, resourceKind); err != nil {
// 		return nil, err
// 	}

// 	return resourceKind, nil
// }

// GetAllEmbeddedFileNames find all resource files from the embedded package.
func (c *Config) GetAllEmbeddedFileNames() ([]string, error) {
	var files []string
	if err := fs.WalkDir(c.eFs, ".", func(filePath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, filePath)

		return err
	}); err != nil {
		return nil, err
	}

	return files, nil
}

// LoadAllEmbeddedResourceFiles load all resource files from the embedded package.
func (c *Config) LoadAllEmbeddedResourceFiles(
	plan bool,
) ([]api.ResourceManager, error) {
	files, err := c.GetAllEmbeddedFileNames()
	if err != nil {
		return nil, err
	}

	fmt.Println(files)
	resources := make([]api.ResourceManager, 0, 1)
	// for _, resourceFileInfo := range files {
	// 	resourceFile, err := loadResourceFile(resourceFileInfo.Path, plan)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	resources = append(resources, resourceFile)
	// }

	return resources, nil
}

// func hashEmbededFile(
// 	eFs embed.FS,
// 	filePath string,
// ) (string, error) {
// 	var returnSHA256String string

// 	file, err := eFs.Open(filePath)
// 	if err != nil {
// 		return returnSHA256String, err
// 	}

// 	defer func() { _ = file.Close() }()

// 	hash := sha256.New()

// 	if _, err := io.Copy(hash, file); err != nil {
// 		return "", err
// 	}
// 	returnSHA256String = hex.EncodeToString(hash.Sum(nil))

// 	return returnSHA256String, nil
// }
