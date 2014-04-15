package cfenv_test

import (
	. "github.com/joefitzgerald/cfenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Application", func() {
	Describe("Application deserialization", func() {
		env := []string{
			`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
			`HOME=/home/vcap/app`,
			`MEMORY_LIMIT=512m`,
			`PWD=/home/vcap`,
			`TMPDIR=/home/vcap/tmp`,
			`USER=vcap`,
			`VCAP_SERVICES={"elephantsql-dev":[{"name":"elephantsql-dev-c6c60","label":"elephantsql-dev","tags":["New Product","relational","Data Store","postgresql"],"plan":"turtle","credentials":{"uri":"postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"}}],"sendgrid":[{"name":"mysendgrid","label":"sendgrid","tags":["smtp","Email"],"plan":"free","credentials":{"hostname":"smtp.sendgrid.net","username":"QvsXMbJ3rK","password":"HCHMOYluTv"}}]}`,
		}
		testEnv := Env(env)
		Context("With valid application variable ", func() {
			It("Should deserialize correctly", func() {
				cfenv := New(testEnv)
				Ω(cfenv).ShouldNot(BeNil())

				Ω(cfenv.ID).Should(BeEquivalentTo("451f045fd16427bb99c895a2649b7b2a"))
				Ω(cfenv.Index).Should(BeEquivalentTo(0))
				Ω(cfenv.Name).Should(BeEquivalentTo("styx-james"))
				Ω(cfenv.Host).Should(BeEquivalentTo("0.0.0.0"))
				Ω(cfenv.Port).Should(BeEquivalentTo(61857))
				Ω(cfenv.Version).Should(BeEquivalentTo("c1063c1c-40b9-434e-a797-db240b587d32"))
				Ω(cfenv.Home).Should(BeEquivalentTo("/home/vcap/app"))
				Ω(cfenv.MemoryLimit).Should(BeEquivalentTo("512m"))
				Ω(cfenv.WorkingDir).Should(BeEquivalentTo("/home/vcap"))
				Ω(cfenv.TempDir).Should(BeEquivalentTo("/home/vcap/tmp"))
				Ω(cfenv.User).Should(BeEquivalentTo("vcap"))
				Ω(len(cfenv.Services)).Should(BeEquivalentTo(2))
				Ω(len(cfenv.Services)).Should(BeEquivalentTo(2))
				Ω(cfenv.Services["elephantsql-dev"][0].Name).Should(BeEquivalentTo("elephantsql-dev-c6c60"))
				Ω(cfenv.Services["elephantsql-dev"][0].Label).Should(BeEquivalentTo("elephantsql-dev"))
				Ω(cfenv.Services["elephantsql-dev"][0].Tags).Should(BeEquivalentTo([]string{"New Product", "relational", "Data Store", "postgresql"}))
				Ω(cfenv.Services["elephantsql-dev"][0].Plan).Should(BeEquivalentTo("turtle"))
				Ω(len(cfenv.Services["elephantsql-dev"][0].Credentials)).Should(BeEquivalentTo(1))
				Ω(cfenv.Services["elephantsql-dev"][0].Credentials["uri"]).Should(BeEquivalentTo("postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"))
				Ω(cfenv.Services["sendgrid"][0].Name).Should(BeEquivalentTo("mysendgrid"))
				Ω(cfenv.Services["sendgrid"][0].Label).Should(BeEquivalentTo("sendgrid"))
				Ω(cfenv.Services["sendgrid"][0].Plan).Should(BeEquivalentTo("free"))
				Ω(cfenv.Services["sendgrid"][0].Tags).Should(BeEquivalentTo([]string{"smtp", "Email"}))
				Ω(len(cfenv.Services["sendgrid"][0].Credentials)).Should(BeEquivalentTo(3))
				Ω(cfenv.Services["sendgrid"][0].Credentials["hostname"]).Should(BeEquivalentTo("smtp.sendgrid.net"))
				Ω(cfenv.Services["sendgrid"][0].Credentials["username"]).Should(BeEquivalentTo("QvsXMbJ3rK"))
				Ω(cfenv.Services["sendgrid"][0].Credentials["password"]).Should(BeEquivalentTo("HCHMOYluTv"))
			})
		})
	})
})
