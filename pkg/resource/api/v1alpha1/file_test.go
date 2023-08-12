package v1alpha1

import (
	// "path/filepath"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	// . "github.com/onsi/gomega"
	// "github.com/spf13/afero"
)

var _ = Describe("File", func() {
	Context("Copy", func() {
		When("dstFile does not exist", func() {
			BeforeEach(func() {
				fileResource := map[string]interface{}{
					"apiVersion": "files.psion.io/v1alpha1",
					"kind":       "File",
					"metadata": map[string]string{
						"name": "file-name",
					},
					"spec": map[string]interface{}{
						"exists": true,
					},
				}

				logger := &logrus.Logger{}
				appFs := afero.NewOsFs()
				f := NewFile(logger, appFs, fileResource)
				fmt.Println(f)
				// f.FileResource.Reconcile(fileResource)
			})

			It("should copy srcFile to dstFile", func() {
				// err := Copy(appFs, srcFile, dstFile)
				// Expect(err).ToNot(HaveOccurred())

				// got := Exists(appFs, dstFile)
				// Expect(got).Should(BeTrue())
			})
		})

		Context("Contents", func() {})
	})
})
