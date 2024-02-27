package file_test

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/file"
)

type FilePublicTestSuite struct {
	suite.Suite

	appDir string
	appFs  afero.Fs
	fm     internal.FileManager
}

func (suite *FilePublicTestSuite) SetupTest() {
	suite.appDir = "/app"
	suite.appFs = afero.NewMemMapFs()
	suite.fm = file.New(suite.appFs)
}

func (suite *FilePublicTestSuite) TestReadOk() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcFile: filepath.Join(suite.appDir, "1.txt"),
		},
	}
	createFileSpecs(specs)

	got, err := suite.fm.Read(specs[0].srcFile)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "mockContent", string(got))
}

func (suite *FilePublicTestSuite) TestReadReturnsErrorWhenFileDoesNotExist() {
	_, err := suite.fm.Read("does-not-exist")
	assert.Error(suite.T(), err)
}

func (suite *FilePublicTestSuite) TestRemoveOk() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcFile: filepath.Join(suite.appDir, "1.txt"),
		},
	}
	createFileSpecs(specs)

	err := suite.fm.Remove(specs[0].srcFile)
	assert.NoError(suite.T(), err)

	got := suite.fm.Exists(specs[0].srcFile)
	assert.False(suite.T(), got)
}

func (suite *FilePublicTestSuite) TestRemoveReturnsErrorWhenFileDoesNotExist() {
	err := suite.fm.Remove("does-not-exist")
	assert.Error(suite.T(), err)
}

func (suite *FilePublicTestSuite) TestExistsOk() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcFile: filepath.Join(suite.appDir, "1.txt"),
		},
	}
	createFileSpecs(specs)

	got := suite.fm.Exists(specs[0].srcFile)
	assert.True(suite.T(), got)
}

func (suite *FilePublicTestSuite) TestExistsReturnsFalseWhenFileDoesNotExist() {
	got := suite.fm.Exists("does-not-exist")
	assert.False(suite.T(), got)
}

func (suite *FilePublicTestSuite) TestGetModeOk() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcFile: filepath.Join(suite.appDir, "1.txt"),
		},
	}
	createFileSpecs(specs)

	got, err := suite.fm.GetMode(specs[0].srcFile)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fs.FileMode(0o644), got)
}

func (suite *FilePublicTestSuite) TestGetModeReturnsErrorWhenFileDoesNotExist() {
	_, err := suite.fm.GetMode("does-not-exist")
	assert.Error(suite.T(), err)
}

func (suite *FilePublicTestSuite) TestGetSetModeOk() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcFile: filepath.Join(suite.appDir, "1.txt"),
		},
	}
	createFileSpecs(specs)

	err := suite.fm.SetMode(specs[0].srcFile, 0o777)
	assert.NoError(suite.T(), err)

	got, _ := suite.fm.GetMode(specs[0].srcFile)
	assert.Equal(suite.T(), fs.FileMode(0o777), got)
}

func (suite *FilePublicTestSuite) TestSetModeReturnsErrorWhenFileDoesNotExist() {
	err := suite.fm.SetMode("does-not-exist", 0o777)
	assert.Error(suite.T(), err)
}

// err := cm.CopyFile(specs[0].srcFile, assertFile)
// got, _ := afero.Exists(suite.appFs, assertFile)
// assert.True(suite.T(), got)

// var _ = Describe("Copy", func() {
// 	appFs := afero.NewMemMapFs()
// 	dir := "/app"
// 	srcFile := filepath.Join(dir, "srcFile")
// 	dstFile := filepath.Join(dir, "dstFile")

// 	BeforeEach(func() {
// 		_ = appFs.MkdirAll(dir, 0o755)
// 	})

// 	When("dstFile does not exist", func() {
// 		BeforeEach(func() {
// 			err := afero.WriteFile(
// 				appFs,
// 				srcFile,
// 				[]byte("mockContent"),
// 				0o644,
// 			)
// 			Expect(err).ToNot(HaveOccurred())
// 		})

// 		It("should copy srcFile to dstFile", func() {
// 			err := Copy(appFs, srcFile, dstFile)
// 			Expect(err).ToNot(HaveOccurred())

// 			got := Exists(appFs, dstFile)
// 			Expect(got).Should(BeTrue())
// 		})
// 	})

// 	When("srcFile does not exist", func() {
// 		It("should have error", func() {
// 			appFs := afero.NewMemMapFs()

// 			err := Copy(appFs, "does-not-exist", "dst")
// 			Expect(err).To(HaveOccurred())
// 		})
// 	})
// })

// var _ = Describe("Size", func() {
// 	When("file exists", func() {
// 		appFs := afero.NewMemMapFs()
// 		dir := "/app"
// 		filePath := filepath.Join(dir, "filePath")

// 		BeforeEach(func() {
// 			_ = appFs.MkdirAll(dir, 0o755)

// 			err := afero.WriteFile(
// 				appFs,
// 				filePath,
// 				[]byte("mockContent"),
// 				0o644,
// 			)
// 			Expect(err).ToNot(HaveOccurred())
// 		})

// 		It("should return file length in bytes", func() {
// 			got, err := Size(appFs, filePath)
// 			Expect(err).ToNot(HaveOccurred())
// 			Expect(got).Should(Equal(int64(11)))
// 		})
// 	})

// 	When("file does not exist", func() {
// 		It("should have error", func() {
// 			appFs := afero.NewMemMapFs()

// 			_, err := Size(appFs, "does-not-exist")
// 			Expect(err).To(HaveOccurred())
// 		})
// 	})
// })

// var _ = Describe("HashFile", func() {
// 	When("file exists", func() {
// 		appFs := afero.NewMemMapFs()
// 		dir := "/app"
// 		filePath := filepath.Join(dir, "filePath")

// 		BeforeEach(func() {
// 			_ = appFs.MkdirAll(dir, 0o755)

// 			err := afero.WriteFile(
// 				appFs,
// 				filePath,
// 				[]byte("mockContent"),
// 				0o644,
// 			)
// 			Expect(err).ToNot(HaveOccurred())
// 		})

// 		It("should return SHA1-hash of file contents", func() {
// 			got, err := HashFile(appFs, filePath)
// 			Expect(err).ToNot(HaveOccurred())
// 			Expect(got).Should(Equal("a388678dad3db361c9198ea665070210e58a0fe5"))
// 		})
// 	})

// 	When("file does not exist", func() {
// 		It("should have error", func() {
// 			appFs := afero.NewMemMapFs()

// 			_, err := HashFile(appFs, "does-not-exist")
// 			Expect(err).To(HaveOccurred())
// 		})
// 	})
// })

// var _ = Describe("Identical", func() {
// 	When("file exists", func() {
// 		appFs := afero.NewMemMapFs()
// 		dir := "/app"
// 		a := filepath.Join(dir, "a")
// 		b := filepath.Join(dir, "b")

// 		BeforeEach(func() {
// 			_ = appFs.MkdirAll(dir, 0o755)

// 			err := afero.WriteFile(
// 				appFs,
// 				a,
// 				[]byte("mockContent"),
// 				0o644,
// 			)
// 			Expect(err).ToNot(HaveOccurred())

// 			err = afero.WriteFile(
// 				appFs,
// 				b,
// 				[]byte("mockContent"),
// 				0o644,
// 			)
// 			Expect(err).ToNot(HaveOccurred())
// 		})

// 		It("should be true", func() {
// 			got, err := Identical(appFs, a, b)
// 			Expect(err).ToNot(HaveOccurred())
// 			Expect(got).Should(BeTrue())
// 		})
// 	})

// 	When("file does not exist", func() {
// 		It("should have error", func() {
// 			appFs := afero.NewMemMapFs()

// 			_, err := Identical(appFs, "does-not-exist-1", "does-not-exist-2")
// 			Expect(err).To(HaveOccurred())
// 		})
// 	})
// })

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestFilePublicTestSuite(t *testing.T) {
	suite.Run(t, new(FilePublicTestSuite))
}
