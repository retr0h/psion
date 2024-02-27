package config_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/config"
	"github.com/retr0h/psion/internal/file"
)

type ConfigPublicTestSuite struct {
	suite.Suite

	appFs afero.Fs
	fm    internal.FileManager
	cm    internal.ConfigManager
}

func (suite *ConfigPublicTestSuite) SetupTest() {
	suite.appFs = afero.NewMemMapFs()
	suite.fm = file.New(suite.appFs)
	suite.cm = config.New(suite.fm)
}

func (suite *ConfigPublicTestSuite) TestLoadConfigOk() {
	runtimeConfigContent := []byte(`
---
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: name
spec:
`)

	err := suite.cm.LoadConfig(runtimeConfigContent)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "name", suite.cm.GetName())
	assert.Equal(suite.T(), "files.psion.io/v1alpha1", suite.cm.GetAPIVersion())
	assert.Equal(suite.T(), "File", suite.cm.GetKind())
}

func (suite *ConfigPublicTestSuite) TestConfigReturnsErrorWhenInvalid() {
	runtimeConfigContent := []byte(`
---
key:
    foo: bar
    path:"bad yaml"
`)
	err := suite.cm.LoadConfig(runtimeConfigContent)
	assert.Error(suite.T(), err)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestConfigPublicTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigPublicTestSuite))
}
