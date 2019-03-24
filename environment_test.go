package cfenv_test

import (
	. "github.com/cloudfoundry-community/go-cfenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
	Describe("Environment variables should be mapped", func() {
		Context("with default environment", func() {
			It("should contain at least one mapped variable", func() {
				vars := CurrentEnv()
				Expect(len(vars)).To(BeNumerically(">", 0), "Environment variables should exist")
			})

			It("should split variables into keys and values", func() {
				vars := CurrentEnv()
				valueCount := 0
				for k, v := range vars {
					// Key should never be empty
					Expect(k).NotTo(BeEmpty())

					// Key should never have equals
					Expect(k).NotTo(ContainSubstring("="))

					// Value may be empty, but let's track non-empty values
					if v != "" {
						valueCount++
					}
				}

				// Ensure we get at least one value from the environment
				Expect(valueCount).To(BeNumerically(">", 0))
			})
		})
	})
})
