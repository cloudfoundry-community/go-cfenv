package cfenv_test

import (
	"os"

	. "github.com/cloudfoundry-community/go-cfenv"
	"github.com/mitchellh/mapstructure"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cfenv", func() {
	Describe("application deserialization", func() {
		validEnv := []string{
			`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"abcabc123123defdef456456","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","space_id":"3e0c28c5-6d9c-436b-b9ee-1f4326e54d05","space_name":"jdk","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
			`HOME=/home/vcap/app`,
			`MEMORY_LIMIT=512m`,
			`PWD=/home/vcap`,
			`TMPDIR=/home/vcap/tmp`,
			`USER=vcap`,
			`VCAP_SERVICES={"elephantsql-dev":[{"name":"elephantsql-dev-c6c60","label":"elephantsql-dev","tags":["New Product","relational","Data Store","postgresql"],"plan":"turtle","credentials":{"uri":"postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"}}],"sendgrid":[{"name":"mysendgrid","label":"sendgrid","tags":["smtp","Email"],"plan":"free","credentials":{"hostname":"smtp.sendgrid.net","username":"QvsXMbJ3rK","password":"HCHMOYluTv"}}],"nfs":[{"credentials":{},"label":"nfs","name":"nfs","plan":"Existing","tags":["nfs"],"volume_mounts":[{"container_dir":"/testpath","device_type":"shared","mode":"rw"}]}]}`,
		}

		validEnvWithoutSpaceIDAndName := []string{
			`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"abcabc123123defdef456456","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
			`HOME=/home/vcap/app`,
			`MEMORY_LIMIT=512m`,
			`PWD=/home/vcap`,
			`TMPDIR=/home/vcap/tmp`,
			`USER=vcap`,
			`VCAP_SERVICES={"elephantsql-dev":[{"name":"elephantsql-dev-c6c60","label":"elephantsql-dev","tags":["New Product","relational","Data Store","postgresql"],"plan":"turtle","credentials":{"uri":"postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"}}],"sendgrid":[{"name":"mysendgrid","label":"sendgrid","tags":["smtp","Email"],"plan":"free","credentials":{"hostname":"smtp.sendgrid.net","username":"QvsXMbJ3rK","password":"HCHMOYluTv"}}]}`,
		}

		envWithIntCredentials := []string{
			`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"abcabc123123defdef456456","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
			`HOME=/home/vcap/app`,
			`MEMORY_LIMIT=512m`,
			`PWD=/home/vcap`,
			`TMPDIR=/home/vcap/tmp`,
			`USER=vcap`,
			`VCAP_SERVICES={"elephantsql-dev":[{"name":"elephantsql-dev-c6c60","label":"elephantsql-dev","tags":["New Product","relational","Data Store","postgresql"],"plan":"turtle","credentials":{"uri":"postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"}}],"cloudantNoSQLDB": [{ "name": "my_cloudant", "label": "cloudantNoSQLDB", "plan": "Shared", "credentials": { "username": "18675309-0000-4aaa-bbbb-999999999-bluemix", "password": "18675309deadbeefaaaabbbbccccddddeeeeffff000099999999999999999999", "host": "01234567-9999-4999-aaaa-abcdefabcdef-bluemix.cloudant.com", "port": 443, "url": "https://18675309-0000-4aaa-bbbb-999999999-bluemix:18675309deadbeefaaaabbbbccccddddeeeeffff000099999999999999999999@01234567-9999-4999-aaaa-abcdefabcdef-bluemix.cloudant.com"}}],"sendgrid":[{"name":"mysendgrid","label":"sendgrid","tags":["smtp","Email"],"plan":"free","credentials":{"hostname":"smtp.sendgrid.net","username":"QvsXMbJ3rK","password":"HCHMOYluTv"}}]}`,
		}

		envWithArrayCredentials := []string{
			`VCAP_APPLICATION={}`,
			`VCAP_SERVICES={"p-kafka": [{"credentials": { "kafka" : { "port": 9092, "node_ips": ["10.244.9.2", "10.244.9.6", "10.244.9.10"]}}}]}`,
		}

		invalidEnv := []string{
			`VCAP_APPLICATION={"instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
			`HOME=/home/vcap/app`,
			`MEMORY_LIMIT_INVALID=512m`,
			`PWD=/home/vcap`,
			`TMPDIR=/home/vcap/tmp`,
			`USER=vcap`,
			`VCAP_SERVICES={"elephantsql-dev":[{"name":"","label":"elephantsql-dev","plan":"turtle","credentials":{"uri":"postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"}}],"sendgrid":[{"name":"mysendgrid","label":"sendgrid","plan":"free","credentials":{"hostname":"smtp.sendgrid.net","username":"QvsXMbJ3rK","password":"HCHMOYluTv"}}]}`,
		}

		notCFEnv := []string{
			`HOME=/home/vcap/app`,
			`MEMORY_LIMIT_INVALID=512m`,
			`PWD=/home/vcap`,
			`PORT=1234`,
			`TMPDIR=/home/vcap/tmp`,
			`USER=vcap`,
		}

		cfEnv := []string{
			`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"abcabc123123defdef456456","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
			`HOME=/home/vcap/app`,
			`MEMORY_LIMIT_INVALID=512m`,
			`PWD=/home/vcap`,
			`PORT=1234`,
			`TMPDIR=/home/vcap/tmp`,
			`USER=vcap`,
			`VCAP_SERVICES={"elephantsql-dev":[{"name":"","label":"elephantsql-dev","plan":"turtle","credentials":{"uri":"postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"}}],"sendgrid":[{"name":"mysendgrid","label":"sendgrid","plan":"free","credentials":{"hostname":"smtp.sendgrid.net","username":"QvsXMbJ3rK","password":"HCHMOYluTv"}}],"nfs":[{"credentials":{},"label":"nfs","name":"nfsexport","plan":"Existing","volume_mounts":[{"container_dir":"/testpath","device_type":"shared","mode":"rw"}]}]}`,
		}

		Context("when not running on Cloud Foundry", func() {
			It("IsRunningOnCF() returns false", func() {
				testEnv := Env(notCFEnv)
				_, err := New(testEnv)
				Expect(err).To(HaveOccurred())
				Expect(IsRunningOnCF()).To(BeFalse())
			})
		})

		Context("when running on Cloud Foundry", func() {
			BeforeEach(func() {
				os.Setenv("VCAP_APPLICATION", "{}")
			})
			AfterEach(func() {
				os.Unsetenv("VCAP_APPLICATION")
			})
			It("IsRunningOnCF() returns true", func() {
				testEnv := Env(cfEnv)
				_, err := New(testEnv)
				Expect(err).NotTo(HaveOccurred())
				Expect(IsRunningOnCF()).To(BeTrue())
			})
		})

		Context("with valid environment", func() {
			It("should deserialize correctly", func() {
				testEnv := Env(validEnv)
				cfenv, err := New(testEnv)
				Expect(err).To(BeNil())
				Expect(cfenv).NotTo(BeNil())

				Expect(cfenv.ID).To(Equal("451f045fd16427bb99c895a2649b7b2a"))
				Expect(cfenv.InstanceID).To(Equal("451f045fd16427bb99c895a2649b7b2a"))
				Expect(cfenv.AppID).To(Equal("abcabc123123defdef456456"))
				Expect(cfenv.CFAPI).To(Equal("https://api.system_domain.com"))
				Expect(cfenv.Index).To(Equal(0))
				Expect(cfenv.Name).To(Equal("styx-james"))
				Expect(cfenv.SpaceName).To(Equal("jdk"))
				Expect(cfenv.SpaceID).To(Equal("3e0c28c5-6d9c-436b-b9ee-1f4326e54d05"))
				Expect(cfenv.Host).To(Equal("0.0.0.0"))
				Expect(cfenv.Port).To(Equal(61857))
				Expect(cfenv.Version).To(Equal("c1063c1c-40b9-434e-a797-db240b587d32"))
				Expect(cfenv.Home).To(Equal("/home/vcap/app"))
				Expect(cfenv.MemoryLimit).To(Equal("512m"))
				Expect(cfenv.WorkingDir).To(Equal("/home/vcap"))
				Expect(cfenv.TempDir).To(Equal("/home/vcap/tmp"))
				Expect(cfenv.User).To(Equal("vcap"))
				Expect(cfenv.Limits.Disk).To(Equal(1024))
				Expect(cfenv.Limits.Mem).To(Equal(512))
				Expect(cfenv.Limits.FDs).To(Equal(16384))
				Expect(cfenv.ApplicationURIs[0]).To(Equal("styx-james.a1-app.cf-app.com"))
				Expect(len(cfenv.Services)).To(Equal(3))
				Expect(cfenv.Services["elephantsql-dev"][0].Name).To(Equal("elephantsql-dev-c6c60"))
				Expect(cfenv.Services["elephantsql-dev"][0].Label).To(Equal("elephantsql-dev"))
				Expect(cfenv.Services["elephantsql-dev"][0].Tags).To(Equal([]string{"New Product", "relational", "Data Store", "postgresql"}))
				Expect(cfenv.Services["elephantsql-dev"][0].Plan).To(Equal("turtle"))
				Expect(len(cfenv.Services["elephantsql-dev"][0].Credentials)).To(Equal(1))
				Expect(cfenv.Services["elephantsql-dev"][0].Credentials["uri"]).To(Equal("postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"))
				Expect(cfenv.Services["sendgrid"][0].Name).To(Equal("mysendgrid"))
				Expect(cfenv.Services["sendgrid"][0].Label).To(Equal("sendgrid"))
				Expect(cfenv.Services["sendgrid"][0].Tags).To(Equal([]string{"smtp", "Email"}))
				Expect(cfenv.Services["sendgrid"][0].Plan).To(Equal("free"))
				Expect(len(cfenv.Services["sendgrid"][0].Credentials)).To(Equal(3))
				Expect(cfenv.Services["sendgrid"][0].Credentials["hostname"]).To(Equal("smtp.sendgrid.net"))
				Expect(cfenv.Services["sendgrid"][0].Credentials["username"]).To(Equal("QvsXMbJ3rK"))
				Expect(cfenv.Services["sendgrid"][0].Credentials["password"]).To(Equal("HCHMOYluTv"))

				Expect(cfenv.Services["nfs"][0].VolumeMounts[0]["container_dir"]).To(Equal("/testpath"))

				name, err := cfenv.Services.WithName("elephantsql-dev-c6c60")
				Expect(name.Name).To(Equal("elephantsql-dev-c6c60"))
				Expect(err).To(BeNil())

				tag, err := cfenv.Services.WithTag("postgresql")
				Expect(len(tag)).To(Equal(1))
				Expect(tag[0].Tags).To(ContainElement("postgresql"))
				Expect(err).To(BeNil())

				label, err := cfenv.Services.WithLabel("elephantsql-dev")
				Expect(len(label)).To(Equal(1))
				Expect(label[0].Label).To(Equal("elephantsql-dev"))
				Expect(err).To(BeNil())

				names, err := cfenv.Services.WithNameUsingPattern(".*(sql|mysend).*")
				Expect(len(names)).To(Equal(2))
				Expect(err).To(BeNil())
				isValidNames := true
				for _, service := range names {
					if service.Name != "mysendgrid" && service.Name != "elephantsql-dev-c6c60" {
						isValidNames = false
					}
				}
				Expect(isValidNames).To(BeTrue(), "Not valid names when finding by regex")

				tags, err := cfenv.Services.WithTagUsingPattern(".*sql.*")
				Expect(len(tags)).To(Equal(1))
				Expect(err).To(BeNil())
				isValidTags := true
				for _, service := range tags {
					if service.Name != "elephantsql-dev-c6c60" {
						isValidTags = false
					}
				}
				Expect(isValidTags).To(BeTrue(), "Not valid tags when finding by regex")

			})

			It("should prefer the PORT environment variable over VCAP_APPLICATION.PORT", func() {
				validEnv = append(validEnv, "PORT=12345")
				testEnv := Env(validEnv)
				cfenv, err := New(testEnv)
				Expect(err).To(BeNil())
				Expect(cfenv).NotTo(BeNil())
				Expect(cfenv.Port).To(Equal(12345))
			})
		})

		Context("without a space name and id", func() {
			It("should deserialize correctly", func() {
				testEnv := Env(validEnvWithoutSpaceIDAndName)
				cfenv, err := New(testEnv)
				Expect(err).To(BeNil())
				Expect(cfenv).NotTo(BeNil())
				Expect(cfenv.SpaceID).To(BeEmpty())
				Expect(cfenv.SpaceName).To(BeEmpty())
			})
		})

		Context("with valid environment with a service with credentials that are an array", func() {
			It("should deserialize correctly", func() {
				testEnv := Env(envWithArrayCredentials)
				cfenv, err := New(testEnv)
				Expect(err).To(BeNil())
				Expect(cfenv).NotTo(BeNil())

				credential := map[string]interface{}{}
				mapstructure.Decode(cfenv.Services["p-kafka"][0].Credentials["kafka"], &credential)

				Expect(len(cfenv.Services["p-kafka"][0].Credentials)).To(Equal(1))
				Expect(credential["node_ips"]).To(Equal([]interface{}{"10.244.9.2", "10.244.9.6", "10.244.9.10"}))
				Expect(credential["port"]).To(Equal(float64(9092)))
			})
		})

		Context("with valid environment with a service with credentials with a port that is an int", func() {
			It("should to deserialize correctly", func() {
				testEnv := Env(envWithIntCredentials)
				cfenv, err := New(testEnv)
				Expect(err).To(BeNil())
				Expect(cfenv).NotTo(BeNil())

				Expect(cfenv.ID).To(Equal("451f045fd16427bb99c895a2649b7b2a"))
				Expect(cfenv.Index).To(Equal(0))
				Expect(cfenv.Name).To(Equal("styx-james"))
				Expect(cfenv.Host).To(Equal("0.0.0.0"))
				Expect(cfenv.Port).To(Equal(61857))
				Expect(cfenv.Version).To(Equal("c1063c1c-40b9-434e-a797-db240b587d32"))
				Expect(cfenv.Home).To(Equal("/home/vcap/app"))
				Expect(cfenv.MemoryLimit).To(Equal("512m"))
				Expect(cfenv.WorkingDir).To(Equal("/home/vcap"))
				Expect(cfenv.TempDir).To(Equal("/home/vcap/tmp"))
				Expect(cfenv.User).To(Equal("vcap"))
				Expect(cfenv.ApplicationURIs[0]).To(Equal("styx-james.a1-app.cf-app.com"))
				Expect(len(cfenv.Services)).To(Equal(3))

				Expect(cfenv.Services["elephantsql-dev"][0].Name).To(Equal("elephantsql-dev-c6c60"))
				Expect(cfenv.Services["elephantsql-dev"][0].Label).To(Equal("elephantsql-dev"))
				Expect(cfenv.Services["elephantsql-dev"][0].Tags).To(Equal([]string{"New Product", "relational", "Data Store", "postgresql"}))
				Expect(cfenv.Services["elephantsql-dev"][0].Plan).To(Equal("turtle"))
				Expect(len(cfenv.Services["elephantsql-dev"][0].Credentials)).To(Equal(1))
				Expect(cfenv.Services["elephantsql-dev"][0].Credentials["uri"]).To(Equal("postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"))

				Expect(cfenv.Services["cloudantNoSQLDB"][0].Name).To(Equal("my_cloudant"))
				Expect(cfenv.Services["cloudantNoSQLDB"][0].Label).To(Equal("cloudantNoSQLDB"))
				Expect(cfenv.Services["cloudantNoSQLDB"][0].Plan).To(Equal("Shared"))
				Expect(len(cfenv.Services["cloudantNoSQLDB"][0].Credentials)).To(Equal(5))
				Expect(cfenv.Services["cloudantNoSQLDB"][0].Credentials["port"]).To(Equal(float64(443)))

				Expect(cfenv.Services["sendgrid"][0].Name).To(Equal("mysendgrid"))
				Expect(cfenv.Services["sendgrid"][0].Label).To(Equal("sendgrid"))
				Expect(cfenv.Services["sendgrid"][0].Tags).To(Equal([]string{"smtp", "Email"}))
				Expect(cfenv.Services["sendgrid"][0].Plan).To(Equal("free"))
				Expect(len(cfenv.Services["sendgrid"][0].Credentials)).To(Equal(3))
				Expect(cfenv.Services["sendgrid"][0].Credentials["hostname"]).To(Equal("smtp.sendgrid.net"))
				Expect(cfenv.Services["sendgrid"][0].Credentials["username"]).To(Equal("QvsXMbJ3rK"))
				Expect(cfenv.Services["sendgrid"][0].Credentials["password"]).To(Equal("HCHMOYluTv"))

				name, err := cfenv.Services.WithName("elephantsql-dev-c6c60")
				Expect(name.Name).To(Equal("elephantsql-dev-c6c60"))
				Expect(err).To(BeNil())

				tag, err := cfenv.Services.WithTag("postgresql")
				Expect(len(tag)).To(Equal(1))
				Expect(tag[0].Tags).To(ContainElement("postgresql"))
				Expect(err).To(BeNil())

				label, err := cfenv.Services.WithLabel("elephantsql-dev")
				Expect(len(label)).To(Equal(1))
				Expect(label[0].Label).To(Equal("elephantsql-dev"))
				Expect(err).To(BeNil())

				names, err := cfenv.Services.WithNameUsingPattern(".*(sql|my_cloud).*")
				Expect(len(names)).To(Equal(2))
				Expect(err).To(BeNil())
				isValidNames := true
				for _, service := range names {
					if service.Name != "my_cloudant" && service.Name != "elephantsql-dev-c6c60" {
						isValidNames = false
					}
				}
				Expect(isValidNames).To(BeTrue(), "Not valid names when finding by regex")

				tags, err := cfenv.Services.WithTagUsingPattern(".*s.*")
				Expect(len(tags)).To(Equal(2))
				Expect(err).To(BeNil())
				isValidTags := true
				for _, service := range tags {
					if service.Name != "mysendgrid" && service.Name != "elephantsql-dev-c6c60" {
						isValidTags = false
					}
				}
				Expect(isValidTags).To(BeTrue(), "Not valid tags when finding by regex")

			})
		})

		Context("with invalid environment", func() {
			It("should deserialize correctly, with missing values", func() {
				testEnv := Env(invalidEnv)
				cfenv, err := New(testEnv)
				Expect(err).To(BeNil())
				Expect(cfenv).NotTo(BeNil())

				Expect(cfenv.ID).To(Equal(""))
				Expect(cfenv.Index).To(Equal(0))
				Expect(cfenv.Name).To(Equal("styx-james"))
				Expect(cfenv.Host).To(Equal("0.0.0.0"))
				Expect(cfenv.Port).To(Equal(61857))
				Expect(cfenv.Version).To(Equal("c1063c1c-40b9-434e-a797-db240b587d32"))
				Expect(cfenv.Home).To(Equal("/home/vcap/app"))
				Expect(cfenv.MemoryLimit).To(Equal(""))
				Expect(cfenv.WorkingDir).To(Equal("/home/vcap"))
				Expect(cfenv.TempDir).To(Equal("/home/vcap/tmp"))
				Expect(cfenv.User).To(Equal("vcap"))
				Expect(cfenv.ApplicationURIs[0]).To(Equal("styx-james.a1-app.cf-app.com"))
				Expect(len(cfenv.Services)).To(Equal(2))
				Expect(len(cfenv.Services)).To(Equal(2))
				Expect(cfenv.Services["elephantsql-dev"][0].Name).To(Equal(""))
				Expect(cfenv.Services["elephantsql-dev"][0].Label).To(Equal("elephantsql-dev"))
				Expect(cfenv.Services["elephantsql-dev"][0].Plan).To(Equal("turtle"))
				Expect(len(cfenv.Services["elephantsql-dev"][0].Credentials)).To(Equal(1))
				Expect(cfenv.Services["elephantsql-dev"][0].Credentials["uri"]).To(Equal("postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"))

				Expect(cfenv.Services["sendgrid"][0].Name).To(Equal("mysendgrid"))
				Expect(cfenv.Services["sendgrid"][0].Label).To(Equal("sendgrid"))
				Expect(cfenv.Services["sendgrid"][0].Plan).To(Equal("free"))
				Expect(len(cfenv.Services["sendgrid"][0].Credentials)).To(Equal(3))
				Expect(cfenv.Services["sendgrid"][0].Credentials["hostname"]).To(Equal("smtp.sendgrid.net"))
				Expect(cfenv.Services["sendgrid"][0].Credentials["username"]).To(Equal("QvsXMbJ3rK"))
				Expect(cfenv.Services["sendgrid"][0].Credentials["password"]).To(Equal("HCHMOYluTv"))
			})
		})
	})

	Describe("CredentialString", func() {
		var service = Service{
			Credentials: map[string]interface{}{
				"string": "stringy-credential",
				"int":    42,
				"nested": map[string]string{
					"key": "value",
				},
			},
		}

		It("returns the requested credential as a string when the credential is a string", func() {
			result, ok := service.CredentialString("string")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal("stringy-credential"))
		})

		It("returns false when the credential is not a string", func() {
			_, ok := service.CredentialString("int")
			Expect(ok).To(BeFalse())
		})

		It("returns false when the credential is a nested thing", func() {
			_, ok := service.CredentialString("nested")
			Expect(ok).To(BeFalse())
		})
	})
})
