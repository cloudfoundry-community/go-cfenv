package cfenv_test

import (
	. "github.com/joefitzgerald/cfenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Application", func() {
	Describe("Application deserialization", func() {
		env := []string{`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","uris":["styx-james.a1-app.cf-app.com"],"users":null}`}
		testEnv := Env(env)
		Context("With valid application variable ", func() {
			It("Should deserialize correctly", func() {
				cfenv := New(testEnv)
				Ω(cfenv).ShouldNot(BeNil())

				Ω(cfenv.Id).Should(BeEquivalentTo("451f045fd16427bb99c895a2649b7b2a"))
				Ω(cfenv.Index).Should(BeEquivalentTo(0))
				Ω(cfenv.Name).Should(BeEquivalentTo("styx-james"))
				Ω(cfenv.Host).Should(BeEquivalentTo("0.0.0.0"))
				Ω(cfenv.Port).Should(BeEquivalentTo(61857))
				Ω(cfenv.Version).Should(BeEquivalentTo("c1063c1c-40b9-434e-a797-db240b587d32"))
			})
		})
	})
})
