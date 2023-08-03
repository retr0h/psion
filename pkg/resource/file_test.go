package resource

import (
	// "path/filepath"
	// "encoding/json"

	"fmt"

	. "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"
	// "github.com/spf13/afero"
)

var _ = Describe("File", func() {
	Context("Copy", func() {})
	Context("Contents", func() {})

	When("dstFile does not exist", func() {
		BeforeEach(func() {
			file := `
        /etc/passwd:
          exists: true
      `
			fmt.Println(file)
		})

		It("should copy srcFile to dstFile", func() {
			// err := Copy(appFs, srcFile, dstFile)
			// Expect(err).ToNot(HaveOccurred())

			// got := Exists(appFs, dstFile)
			// Expect(got).Should(BeTrue())
		})
	})
})
