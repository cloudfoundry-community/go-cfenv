package cfenv_test

import (
	. "github.com/joefitzgerald/cfenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Envmap", func() {
	Describe("Environment variables should be mapped", func() {
		Context("With default environment", func() {
			It("Should contain at least one mapped variable", func() {
				vars := Variables()
				Ω(len(vars)).Should(BeNumerically(">", 0), "Environment variables should exist")
			})

			It("Should split variables into keys and values", func() {
				vars := Variables()
				for k, v := range vars {
					Ω(k).ShouldNot(BeEmpty())
					Ω(v).ShouldNot(BeEmpty())
				}
			})
		})
	})
})
