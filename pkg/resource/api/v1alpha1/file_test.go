package v1alpha1

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"

	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
)

var testFileResourceContent []byte

var _ = Describe("File", func() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
	Context("plan", func() {
		plan := true

		When("file exists", func() {
			appFs := afero.NewMemMapFs()

			BeforeEach(func() {
				_ = appFs.MkdirAll(dir, 0o755)

				err := afero.WriteFile(
					appFs,
					filePath,
					[]byte("mockContent"),
					0o644,
				)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})

			It("should plan to remove file", func() {
				var resource api.Manager = NewFile(
					logger,
					appFs,
					plan,
				)

				err := yaml.Unmarshal(testFileResourceContent, resource)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				err = resource.Reconcile()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				got := file.Exists(appFs, filePath)
				gomega.Expect(got).Should(gomega.BeTrue())

				gomega.Expect(resource.GetStatus()).To(gomega.Equal(api.Pending))
				gomega.Expect(resource.GetMessage()).To(gomega.Equal("file exists"))
				gomega.Expect(resource.GetReason()).To(gomega.Equal("Plan"))
			})
		})

		When("file does not exist", func() {
			appFs := afero.NewMemMapFs()

			BeforeEach(func() {
				_ = appFs.MkdirAll(dir, 0o755)
			})

			It("should plan not remove file", func() {
				var resource api.Manager = NewFile(
					logger,
					appFs,
					plan,
				)

				err := yaml.Unmarshal(testFileResourceContent, resource)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				err = resource.Reconcile()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				gomega.Expect(resource.GetStatus()).To(gomega.Equal(api.Pending))
				gomega.Expect(resource.GetMessage()).To(gomega.Equal("file does not exist"))
				gomega.Expect(resource.GetReason()).To(gomega.Equal("Plan"))
			})
		})
	})
})
