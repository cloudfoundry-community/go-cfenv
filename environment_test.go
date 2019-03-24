package cfenv_test

import (
	"testing"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testEnvironment(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("environment variables should be mapped with default environment", func() {
		it("should contain at least one mapped variable", func() {
			vars := cfenv.CurrentEnv()
			Expect(len(vars)).To(BeNumerically(">", 0), "Environment variables should exist")
		})

		it("should split variables into keys and values", func() {
			vars := cfenv.CurrentEnv()
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
}
