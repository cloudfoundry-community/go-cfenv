package cfenv

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Envmap", func() {
	Describe("Environment variables should be split correctly", func() {
		test := func(input string, expectedKey string, expectedValue string) {
			k, v := splitEnv(input)
			Expect(k).To(Equal(expectedKey))
			Expect(v).To(Equal(expectedValue))
		}

		Context("with empty env var", func() {
			It("should have empty value", func() {
				test("", "", "")
			})
		})

		Context("with env var not split by equals", func() {
			It("should have empty value", func() {
				test("TEST", "TEST", "")
			})
		})

		Context("with env var split by equals but no value", func() {
			It("should have empty value", func() {
				test("TEST=", "TEST", "")
			})
		})

		Context("with env var split by equals with key and value", func() {
			It("should have non-empty key and value", func() {
				test("TEST=VAL", "TEST", "VAL")
			})
		})

		Context("with env var split by equals with key and value containing equals", func() {
			It("should have non-empty key and value", func() {
				test("TEST=VAL=OTHERVAL", "TEST", "VAL=OTHERVAL")
			})
		})
	})
})
