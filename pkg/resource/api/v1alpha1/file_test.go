package v1alpha1

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
				Expect(err).ToNot(HaveOccurred())
			})

			It("should plan to remove file", func() {
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
				Expect(got).Should(BeTrue())

				Expect(resource.GetStatus()).To(Equal(api.Pending))

				conditions := resource.GetStatusConditions()
				Expect(conditions).To(HaveLen(1))
				Expect(conditions[0].GetType()).To(Equal("Remove"))
				Expect(conditions[0].GetStatus()).To(Equal(api.Pending))
				Expect(conditions[0].GetMessage()).To(Equal("file exists"))
				Expect(conditions[0].GetReason()).To(Equal(Plan))
				Expect(conditions[0].GetGot()).To(Equal("file exists"))
				Expect(conditions[0].GetWant()).To(Equal(NoOp))
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
				Expect(err).ToNot(HaveOccurred())

				err = resource.Reconcile()
				Expect(err).ToNot(HaveOccurred())

				Expect(resource.GetStatus()).To(Equal(api.Pending))

				conditions := resource.GetStatusConditions()
				Expect(conditions).To(HaveLen(1))
				Expect(conditions[0].GetType()).To(Equal("Remove"))
				Expect(conditions[0].GetStatus()).To(Equal(api.Pending))
				Expect(conditions[0].GetMessage()).To(Equal("file does not exist"))
				Expect(conditions[0].GetReason()).To(Equal(Plan))
				Expect(conditions[0].GetGot()).To(Equal("file does not exist"))
				Expect(conditions[0].GetWant()).To(Equal(NoOp))
			})
		})
	})

	Context("apply", func() {
		plan := false

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

				conditions := resource.GetStatusConditions()
				Expect(conditions).To(HaveLen(1))
				Expect(conditions[0].GetType()).To(Equal("Remove"))
				Expect(conditions[0].GetStatus()).To(Equal(api.Succeeded))
				Expect(conditions[0].GetMessage()).To(Equal("file removed"))
				Expect(conditions[0].GetReason()).To(Equal(Apply))
				Expect(conditions[0].GetGot()).To(Equal("file exists"))
				Expect(conditions[0].GetWant()).To(Equal("file removed"))
			})
		})

		When("file does not exist", func() {
			appFs := afero.NewMemMapFs()

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

				Expect(resource.GetStatus()).To(Equal(api.Succeeded))

				conditions := resource.GetStatusConditions()
				Expect(conditions).To(HaveLen(1))
				Expect(conditions[0].GetType()).To(Equal("Remove"))
				Expect(conditions[0].GetStatus()).To(Equal(api.Succeeded))
				Expect(conditions[0].GetMessage()).To(Equal("file does not exist"))
				Expect(conditions[0].GetReason()).To(Equal(Apply))
				Expect(conditions[0].GetGot()).To(Equal("file does not exist"))
				Expect(conditions[0].GetWant()).To(Equal(NoOp))
			})
		})
	})
})
