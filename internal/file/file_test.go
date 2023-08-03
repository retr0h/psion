package file

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Copy", func() {
	appFs := afero.NewMemMapFs()
	srcDir := "/src"
	dstDir := "/dst"
	srcFile := filepath.Join(srcDir, "srcFile")
	dstFile := filepath.Join(dstDir, "dstFile")

	BeforeEach(func() {
		_ = appFs.MkdirAll(srcDir, 0o755)
		_ = appFs.MkdirAll(dstDir, 0o755)
	})

	When("dstFile does not exist", func() {
		BeforeEach(func() {
			err := afero.WriteFile(
				appFs,
				srcFile,
				[]byte("mockContent"),
				0o644,
			)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should copy srcFile to dstFile", func() {
			err := Copy(appFs, srcFile, dstFile)
			Expect(err).ToNot(HaveOccurred())

			got := Exists(appFs, dstFile)
			Expect(got).Should(BeTrue())
		})
	})

	When("srcFile does not exist", func() {
		It("should have error", func() {
			appFs := afero.NewMemMapFs()

			err := Copy(appFs, "does-not-exist", "dst")
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Exists", func() {
	When("srcFile exists", func() {
		appFs := afero.NewMemMapFs()
		srcDir := "/src"
		srcFile := filepath.Join(srcDir, "srcFile")

		BeforeEach(func() {
			_ = appFs.MkdirAll(srcDir, 0o755)

			err := afero.WriteFile(
				appFs,
				srcFile,
				[]byte("mockContent"),
				0o644,
			)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should be true", func() {
			got := Exists(appFs, srcFile)
			Expect(got).Should(BeTrue())
		})
	})

	When("srcFile does not exist", func() {
		It("should be false", func() {
			appFs := afero.NewMemMapFs()

			got := Exists(appFs, "does-not-exist")
			Expect(got).Should(BeFalse())
		})
	})
})

var _ = Describe("Size", func() {
	When("file exists", func() {
		appFs := afero.NewMemMapFs()
		srcDir := "/src"
		srcFile := filepath.Join(srcDir, "srcFile")

		BeforeEach(func() {
			_ = appFs.MkdirAll(srcDir, 0o755)

			err := afero.WriteFile(
				appFs,
				srcFile,
				[]byte("mockContent"),
				0o644,
			)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return file length in bytes", func() {
			got, err := Size(appFs, srcFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(got).Should(Equal(int64(11)))
		})
	})

	When("file does not exist", func() {
		It("should have error", func() {
			appFs := afero.NewMemMapFs()

			_, err := Size(appFs, "does-not-exist")
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("HashFile", func() {
	When("file exists", func() {
		appFs := afero.NewMemMapFs()
		srcDir := "/src"
		srcFile := filepath.Join(srcDir, "srcFile")

		BeforeEach(func() {
			_ = appFs.MkdirAll(srcDir, 0o755)

			err := afero.WriteFile(
				appFs,
				srcFile,
				[]byte("mockContent"),
				0o644,
			)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return SHA1-hash of file contents", func() {
			got, err := HashFile(appFs, srcFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(got).Should(Equal("a388678dad3db361c9198ea665070210e58a0fe5"))
		})
	})

	When("file does not exist", func() {
		It("should have error", func() {
			appFs := afero.NewMemMapFs()

			_, err := HashFile(appFs, "does-not-exist")
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Identical", func() {
	When("file exists", func() {
		appFs := afero.NewMemMapFs()
		srcDir := "/src"
		dstDir := "/dst"
		srcFile := filepath.Join(srcDir, "srcFile")
		dstFile := filepath.Join(dstDir, "dstFile")

		BeforeEach(func() {
			_ = appFs.MkdirAll(srcDir, 0o755)

			err := afero.WriteFile(
				appFs,
				srcFile,
				[]byte("mockContent"),
				0o644,
			)
			Expect(err).ToNot(HaveOccurred())

			err = afero.WriteFile(
				appFs,
				dstFile,
				[]byte("mockContent"),
				0o644,
			)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should be true", func() {
			got, err := Identical(appFs, srcFile, dstFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(got).Should(BeTrue())
		})
	})

	When("file does not exist", func() {
		It("should have error", func() {
			appFs := afero.NewMemMapFs()

			_, err := Identical(appFs, "does-not-exist-1", "does-not-exist-2")
			Expect(err).To(HaveOccurred())
		})
	})
})
