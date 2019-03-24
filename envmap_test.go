package cfenv

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testEnvMap(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	test := func(input string, expectedKey string, expectedValue string) {
		k, v := splitEnv(input)
		Expect(k).To(Equal(expectedKey))
		Expect(v).To(Equal(expectedValue))
	}

	when("splitting environment variables", func() {
		when("with empty env var", func() {
			it("should have empty value", func() {
				test("", "", "")
			})
		})

		when("with env var not split by equals", func() {
			it("should have empty value", func() {
				test("TEST", "TEST", "")
			})
		})

		when("with env var split by equals but no value", func() {
			it("should have empty value", func() {
				test("TEST=", "TEST", "")
			})
		})

		when("with env var split by equals with key and value", func() {
			it("should have non-empty key and value", func() {
				test("TEST=VAL", "TEST", "VAL")
			})
		})

		when("with env var split by equals with key and value containing equals", func() {
			it("should have non-empty key and value", func() {
				test("TEST=VAL=OTHERVAL", "TEST", "VAL=OTHERVAL")
			})
		})
	})
}
