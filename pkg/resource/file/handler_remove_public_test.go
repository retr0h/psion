package file_test

import (
	"fmt"
	// "io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/yaml"

	"github.com/retr0h/psion/internal"
	intFile "github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource"
	"github.com/retr0h/psion/pkg/resource/api"
	"github.com/retr0h/psion/pkg/resource/file"
)

type RemoveHandlerPublicTestSuite struct {
	suite.Suite

	appDir          string
	filePath        string
	appFs           afero.Fs
	logger          *slog.Logger
	f               internal.FileManager
	resourceContent []byte
}

func (suite *RemoveHandlerPublicTestSuite) SetupTest() {
	suite.appDir = "/app"
	suite.filePath = filepath.Join(suite.appDir, "filePath")
	suite.appFs = afero.NewMemMapFs()
	suite.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	suite.f = intFile.New(suite.appFs)
	suite.resourceContent = []byte(fmt.Sprintf(`
---
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: name
spec:
  exists: false
  path: %s
`, suite.filePath))

	_ = suite.appFs.MkdirAll(suite.appDir, 0o755)
}

func (suite *RemoveHandlerPublicTestSuite) TestPlanReconcileRemoveOk() {
	_ = afero.WriteFile(
		suite.appFs,
		suite.filePath,
		[]byte("mockContent"),
		0o644,
	)

	plan := true
	var resource resource.Manager = file.New(
		suite.logger,
		suite.f,
		plan,
	)

	err := yaml.Unmarshal(suite.resourceContent, resource)
	assert.NoError(suite.T(), err)

	err = resource.Reconcile()
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), api.Pending, resource.GetStatus())

	conditions := resource.GetStatusConditions()
	assert.Len(suite.T(), conditions, 1)
	assert.Equal(suite.T(), file.RemoveAction, conditions[0].GetType())
	assert.Equal(suite.T(), api.Pending, conditions[0].GetStatus())
	assert.Equal(suite.T(), "file exists", conditions[0].GetMessage())
	assert.Equal(suite.T(), api.Plan, conditions[0].GetReason())
	assert.Equal(suite.T(), "exists true", conditions[0].GetGot())
	assert.Equal(suite.T(), "exists false", conditions[0].GetWant())
}

func (suite *RemoveHandlerPublicTestSuite) TestPlanReconcileRemoveOkWhenFileDoesNotExist() {
	plan := true
	var resource resource.Manager = file.New(
		suite.logger,
		suite.f,
		plan,
	)

	err := yaml.Unmarshal(suite.resourceContent, resource)
	assert.NoError(suite.T(), err)

	err = resource.Reconcile()
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), api.NoOp, resource.GetStatus())

	conditions := resource.GetStatusConditions()
	assert.Len(suite.T(), conditions, 1)
	assert.Equal(suite.T(), file.RemoveAction, conditions[0].GetType())
	assert.Equal(suite.T(), api.NoOp, conditions[0].GetStatus())
	assert.Equal(suite.T(), "file does not exist", conditions[0].GetMessage())
	assert.Equal(suite.T(), api.Plan, conditions[0].GetReason())
	assert.Equal(suite.T(), "exists false", conditions[0].GetGot())
	assert.Equal(suite.T(), "exists false", conditions[0].GetWant())
}

func (suite *RemoveHandlerPublicTestSuite) TestApplyReconcileRemoveOk() {
	_ = afero.WriteFile(
		suite.appFs,
		suite.filePath,
		[]byte("mockContent"),
		0o644,
	)

	plan := false
	var resource resource.Manager = file.New(
		suite.logger,
		suite.f,
		plan,
	)

	err := yaml.Unmarshal(suite.resourceContent, resource)
	assert.NoError(suite.T(), err)

	err = resource.Reconcile()
	assert.NoError(suite.T(), err)

	got := suite.f.Exists(suite.filePath)
	assert.False(suite.T(), got)

	assert.Equal(suite.T(), api.Succeeded, resource.GetStatus())

	conditions := resource.GetStatusConditions()
	assert.Len(suite.T(), conditions, 1)
	assert.Equal(suite.T(), file.RemoveAction, conditions[0].GetType())
	assert.Equal(suite.T(), api.Succeeded, conditions[0].GetStatus())
	assert.Equal(suite.T(), "file removed", conditions[0].GetMessage())
	assert.Equal(suite.T(), api.Apply, conditions[0].GetReason())
	assert.Equal(suite.T(), "exists false", conditions[0].GetGot())
	assert.Equal(suite.T(), "exists false", conditions[0].GetWant())
}

func (suite *RemoveHandlerPublicTestSuite) TestApplyReconcileRemoveOkWhenFileDoesNotExist() {
	plan := false
	var resource resource.Manager = file.New(
		suite.logger,
		suite.f,
		plan,
	)

	err := yaml.Unmarshal(suite.resourceContent, resource)
	assert.NoError(suite.T(), err)

	err = resource.Reconcile()
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), api.NoOp, resource.GetStatus())

	conditions := resource.GetStatusConditions()
	assert.Len(suite.T(), conditions, 1)
	assert.Equal(suite.T(), file.RemoveAction, conditions[0].GetType())
	assert.Equal(suite.T(), api.NoOp, conditions[0].GetStatus())
	assert.Equal(suite.T(), "file does not exist", conditions[0].GetMessage())
	assert.Equal(suite.T(), api.Apply, conditions[0].GetReason())
	assert.Equal(suite.T(), "exists false", conditions[0].GetGot())
	assert.Equal(suite.T(), "exists false", conditions[0].GetWant())
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRemoveHanderPublicTestSuite(t *testing.T) {
	suite.Run(t, new(RemoveHandlerPublicTestSuite))
}
