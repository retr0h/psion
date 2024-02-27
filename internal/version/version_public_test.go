package version_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/version"
)

type VersionPublicTestSuite struct {
	suite.Suite

	vm internal.VersionManager

	ver    string
	commit string
	date   string
}

func (suite *VersionPublicTestSuite) SetupTest() {
	suite.vm = version.New()

	suite.ver = "dev"
	suite.commit = "none"
	suite.date = "unknown"
}

func (suite *VersionPublicTestSuite) TestLoadVersionToShortenedOk() {
	suite.vm.LoadVersion(suite.ver, suite.commit, suite.date)
	got := suite.vm.ToShortened()
	want := fmt.Sprintf("Version: %s\n", suite.ver)

	assert.Equal(suite.T(), want, got)
}

func (suite *VersionPublicTestSuite) TestLoadVersionToJSONOk() {
	suite.vm.LoadVersion(suite.ver, suite.commit, suite.date)
	got := suite.vm.ToJSON()
	want, err := json.Marshal(
		struct {
			Version string `json:"version"`
			Commit  string `json:"commit"`
			Date    string `json:"date"`
		}{
			suite.ver,
			suite.commit,
			suite.date,
		})

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), string(want), got)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestVersionPublicTestSuite(t *testing.T) {
	suite.Run(t, new(VersionPublicTestSuite))
}
