package v1alpha1

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"

	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
)

var testFileResourceContent []byte

var _ = Describe("File", func() {
	logger := &logrus.Logger{}
	appFs := afero.NewMemMapFs()
	plan := false

	dir := "/app"
	filePath := filepath.Join(dir, "filePath")
	testFileResourceContent = []byte(fmt.Sprintf(`
---
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: name
spec:
  exists: false
  path: %s
`, filePath))
	Context("apply", func() {
		When("file exists", func() {
			BeforeEach(func() {
				_ = appFs.MkdirAll(dir, 0o755)

				err := afero.WriteFile(
					appFs,
					filePath,
					[]byte("mockContent"),
					0o644,
				)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should remove file", func() {
				var resource api.Manager = NewFile(
					logger,
					appFs,
					plan,
				)

				err := yaml.Unmarshal(testFileResourceContent, resource)
				Expect(err).ToNot(HaveOccurred())

				err = resource.Reconcile()
				Expect(err).ToNot(HaveOccurred())

				got := file.Exists(appFs, filePath)
				Expect(got).Should(BeFalse())

				Expect(resource.GetStatus()).To(Equal(api.Succeeded))
				Expect(resource.GetMessage()).To(Equal("remove: /app/filePath"))
				Expect(resource.GetReason()).To(Equal(""))
			})
		})

		When("file does not exist", func() {
			logger := &logrus.Logger{}

			BeforeEach(func() {
				_ = appFs.MkdirAll(dir, 0o755)
			})

			It("should not remove file", func() {
				var resource api.Manager = NewFile(
					logger,
					appFs,
					plan,
				)

				err := yaml.Unmarshal(testFileResourceContent, resource)
				Expect(err).ToNot(HaveOccurred())

				err = resource.Reconcile()
				Expect(err).ToNot(HaveOccurred())

				Expect(resource.GetStatus()).To(Equal(api.Failed))
				Expect(resource.GetMessage()).To(Equal("remove: /app/filePath"))
				Expect(resource.GetReason()).To(Equal("does not exist"))
			})
		})
	})
})
