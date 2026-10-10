package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/config"
)

type ConfigPublicTestSuite struct {
	suite.Suite

	c internal.ConfigManager
}

func (suite *ConfigPublicTestSuite) SetupTest() {
	suite.c = config.New()
}

func (suite *ConfigPublicTestSuite) TestConfigOk() {
	runtimeConfigContent := []byte(`
---
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: name
spec:
`)

	got, err := suite.c.GetConfig(runtimeConfigContent)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "name", got.Name)
	assert.Equal(suite.T(), "files.psion.io/v1alpha1", got.APIVersion)
	assert.Equal(suite.T(), "File", got.Kind)
}

func (suite *ConfigPublicTestSuite) TestConfigReturnsErrorWhenInvalid() {
	runtimeConfigContent := []byte(`
---
key:
    foo: bar
    path:"bad yaml"
`)
	_, err := suite.c.GetConfig(runtimeConfigContent)
	assert.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, config.ErrInvalidConfig)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestConfigPublicTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigPublicTestSuite))
}
