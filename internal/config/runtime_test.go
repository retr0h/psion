package config

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadRuntimeConfig", func() {
	var runtimeConfigContent []byte

	When("config is valid", func() {
		BeforeEach(func() {
			runtimeConfigContent = []byte(`
---
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: name
spec:
`)
		})

		It("should return config", func() {
			got, err := LoadRuntimeConfig(runtimeConfigContent)
			Expect(err).ToNot(HaveOccurred())

			Expect(got.Name).Should(Equal("name"))
			Expect(got.APIVersion).Should(Equal("files.psion.io/v1alpha1"))
			Expect(got.Kind).Should(Equal("File"))
		})
	})

	When("config is invalid", func() {
		BeforeEach(func() {
			runtimeConfigContent = []byte(`
---
key:
    foo: bar
    path:"bad yaml"
`)
		})

		It("should return an error", func() {
			_, err := LoadRuntimeConfig(runtimeConfigContent)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ErrInvalidRuntimeConfig))
		})
	})
})
