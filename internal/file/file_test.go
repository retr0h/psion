package file

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Read", func() {
	When("file exists", func() {
		appFs := afero.NewMemMapFs()
		dir := "/app"
		filePath := filepath.Join(dir, "filePath")

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

		It("should return file []byte", func() {
			got, err := Read(appFs, filePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(got)).Should(Equal("mockContent"))

			got, err = Read(appFs, filePath)
			fmt.Println(got)
			fmt.Println(err)
		})
	})

	When("file does not exist", func() {
		It("should have error", func() {
			appFs := afero.NewMemMapFs()

			_, err := Read(appFs, "does-not-exist")
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Copy", func() {
	appFs := afero.NewMemMapFs()
	dir := "/app"
	srcFile := filepath.Join(dir, "srcFile")
	dstFile := filepath.Join(dir, "dstFile")

	BeforeEach(func() {
		_ = appFs.MkdirAll(dir, 0o755)
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
	When("file exists", func() {
		appFs := afero.NewMemMapFs()
		dir := "/app"
		filePath := filepath.Join(dir, "filePath")

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

		It("should be true", func() {
			got := Exists(appFs, filePath)
			Expect(got).Should(BeTrue())
		})
	})

	When("file does not exist", func() {
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
		dir := "/app"
		filePath := filepath.Join(dir, "filePath")

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

		It("should return file length in bytes", func() {
			got, err := Size(appFs, filePath)
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
		dir := "/app"
		filePath := filepath.Join(dir, "filePath")

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

		It("should return SHA1-hash of file contents", func() {
			got, err := HashFile(appFs, filePath)
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
		dir := "/app"
		a := filepath.Join(dir, "a")
		b := filepath.Join(dir, "b")

		BeforeEach(func() {
			_ = appFs.MkdirAll(dir, 0o755)

			err := afero.WriteFile(
				appFs,
				a,
				[]byte("mockContent"),
				0o644,
			)
			Expect(err).ToNot(HaveOccurred())

			err = afero.WriteFile(
				appFs,
				b,
				[]byte("mockContent"),
				0o644,
			)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should be true", func() {
			got, err := Identical(appFs, a, b)
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
