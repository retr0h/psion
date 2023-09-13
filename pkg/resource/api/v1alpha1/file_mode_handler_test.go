package v1alpha1

import (
	"fmt"
	"io/fs"
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

var _ = Describe("file handler", func() {
	var testFileResourceContent []byte
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
  exists: true
  path: %s
  mode: 0o700
`, filePath))
	Context("plan", func() {
		plan := true

		When("file exists", func() {
			When("modes differ", func() {
				appFs := afero.NewMemMapFs()
				fileManager := file.New(appFs)

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

				It("should plan to update file modes", func() {
					var resource api.Manager = NewFile(
						logger,
						fileManager,
						plan,
					)

					err := yaml.Unmarshal(testFileResourceContent, resource)
					Expect(err).ToNot(HaveOccurred())

					err = resource.Reconcile()
					Expect(err).ToNot(HaveOccurred())

					got := fileManager.Exists(filePath)
					Expect(got).Should(BeTrue())

					Expect(resource.GetStatus()).To(Equal(api.Pending))

					conditions := resource.GetStatusConditions()
					Expect(conditions).To(HaveLen(1))
					Expect(conditions[0].GetType()).To(Equal(ModeAction))
					Expect(conditions[0].GetStatus()).To(Equal(api.Pending))
					Expect(conditions[0].GetMessage()).To(Equal("modes differ"))
					Expect(conditions[0].GetReason()).To(Equal(api.Plan))
					Expect(conditions[0].GetGot()).To(Equal("0644"))
					Expect(conditions[0].GetWant()).To(Equal("0700"))
				})
			})

			When("modes are the same", func() {
				appFs := afero.NewMemMapFs()
				fileManager := file.New(appFs)

				BeforeEach(func() {
					_ = appFs.MkdirAll(dir, 0o755)

					err := afero.WriteFile(
						appFs,
						filePath,
						[]byte("mockContent"),
						0o700,
					)
					Expect(err).ToNot(HaveOccurred())
				})

				It("should not plan to update file modes", func() {
					var resource api.Manager = NewFile(
						logger,
						fileManager,
						plan,
					)

					err := yaml.Unmarshal(testFileResourceContent, resource)
					Expect(err).ToNot(HaveOccurred())

					err = resource.Reconcile()
					Expect(err).ToNot(HaveOccurred())

					got := fileManager.Exists(filePath)
					Expect(got).Should(BeTrue())

					Expect(resource.GetStatus()).To(Equal(api.NoOp))

					conditions := resource.GetStatusConditions()
					Expect(conditions).To(HaveLen(1))
					Expect(conditions[0].GetType()).To(Equal(ModeAction))
					Expect(conditions[0].GetStatus()).To(Equal(api.NoOp))
					Expect(conditions[0].GetMessage()).To(Equal("modes same"))
					Expect(conditions[0].GetReason()).To(Equal(api.Plan))
					Expect(conditions[0].GetGot()).To(Equal("0700"))
					Expect(conditions[0].GetWant()).To(Equal("0700"))
				})
			})
		})

		When("chown fails", func() {
			appFs := afero.NewMemMapFs()
			fileManager := file.New(appFs)

			BeforeEach(func() {
				_ = appFs.MkdirAll(dir, 0o755)
			})

			It("should have error", func() {
				var resource api.Manager = NewFile(
					logger,
					fileManager,
					plan,
				)

				err := yaml.Unmarshal(testFileResourceContent, resource)
				Expect(err).ToNot(HaveOccurred())

				err = resource.Reconcile()
				Expect(err).ToNot(HaveOccurred())

				Expect(resource.GetStatus()).To(Equal(api.Failed))

				conditions := resource.GetStatusConditions()
				Expect(conditions).To(HaveLen(1))
				Expect(conditions[0].GetType()).To(Equal(ModeAction))
				Expect(conditions[0].GetStatus()).To(Equal(api.Failed))
				Expect(
					conditions[0].GetMessage(),
				).To(Equal("open /app/filePath: file does not exist"))
				Expect(conditions[0].GetReason()).To(Equal(api.Plan))
				Expect(conditions[0].GetGot()).To(Equal("Unknown"))
				Expect(conditions[0].GetWant()).To(Equal("0700"))
			})
		})
	})

	Context("apply", func() {
		plan := false

		When("file exists", func() {
			When("modes differ", func() {
				appFs := afero.NewMemMapFs()
				fileManager := file.New(appFs)

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

				It("should update file modes", func() {
					var resource api.Manager = NewFile(
						logger,
						fileManager,
						plan,
					)

					err := yaml.Unmarshal(testFileResourceContent, resource)
					Expect(err).ToNot(HaveOccurred())

					err = resource.Reconcile()
					Expect(err).ToNot(HaveOccurred())

					got, err := fileManager.GetMode(filePath)
					Expect(err).ToNot(HaveOccurred())
					Expect(got).Should(Equal(fs.FileMode(0o700)))

					Expect(resource.GetStatus()).To(Equal(api.Succeeded))

					conditions := resource.GetStatusConditions()
					Expect(conditions).To(HaveLen(1))
					Expect(conditions[0].GetType()).To(Equal(ModeAction))
					Expect(conditions[0].GetStatus()).To(Equal(api.Succeeded))
					Expect(conditions[0].GetMessage()).To(Equal("modes updated"))
					Expect(conditions[0].GetReason()).To(Equal(api.Apply))
					Expect(conditions[0].GetGot()).To(Equal("0700"))
					Expect(conditions[0].GetWant()).To(Equal("0700"))
				})
			})

			When("set mode fails", func() {
				BeforeEach(func() {
				})

				It("should have error", func() {
				})
			})

			When("modes are the same", func() {
				appFs := afero.NewMemMapFs()
				fileManager := file.New(appFs)

				BeforeEach(func() {
					_ = appFs.MkdirAll(dir, 0o755)

					err := afero.WriteFile(
						appFs,
						filePath,
						[]byte("mockContent"),
						0o700,
					)
					Expect(err).ToNot(HaveOccurred())
				})

				It("should not update file modes", func() {
					var resource api.Manager = NewFile(
						logger,
						fileManager,
						plan,
					)

					err := yaml.Unmarshal(testFileResourceContent, resource)
					Expect(err).ToNot(HaveOccurred())

					err = resource.Reconcile()
					Expect(err).ToNot(HaveOccurred())

					Expect(resource.GetStatus()).To(Equal(api.NoOp))

					conditions := resource.GetStatusConditions()
					Expect(conditions).To(HaveLen(1))
					Expect(conditions[0].GetType()).To(Equal(ModeAction))
					Expect(conditions[0].GetStatus()).To(Equal(api.NoOp))
					Expect(conditions[0].GetMessage()).To(Equal("modes same"))
					Expect(conditions[0].GetReason()).To(Equal(api.Apply))
					Expect(conditions[0].GetGot()).To(Equal("0700"))
					Expect(conditions[0].GetWant()).To(Equal("0700"))
				})
			})
		})
	})
})
