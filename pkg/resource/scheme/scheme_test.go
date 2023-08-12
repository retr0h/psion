package scheme

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Scheme", func() {
	Context("Version", func() {
		var scheme *Scheme

		When("apiVersion is correct", func() {
			BeforeEach(func() {
				logger := &logrus.Logger{}
				scheme = New(logger)
			})
		})

		It("should have Version", func() {
			resourceYAML := []byte(`
---
apiVersion: files.psion.io/v1alpha1
kind: File
`)
			got := scheme.Decode(resourceYAML).Version()
			Expect(got).Should(Equal("v1alpha1"))
		})

		When("apiVersion cannot parsed", func() {
			BeforeEach(func() {
				logger := &logrus.Logger{}
				scheme = New(logger)
			})
		})

		It("should not have Version", func() {
			resourceYAML := []byte(`
---
apiVersion: files.psion.io
kind: File
`)
			got := scheme.Decode(resourceYAML).Version()
			Expect(got).Should(Equal(""))
		})
	})

	Context("Kind", func() {
		var scheme *Scheme

		It("should have Kind", func() {
			logger := &logrus.Logger{}
			scheme = New(logger)
			resourceYAML := []byte(`
---
apiVersion: files.psion.io/v1alpha1
kind: File
`)
			got := scheme.Decode(resourceYAML).Kind
			Expect(got).Should(Equal("File"))
		})
	})
})
