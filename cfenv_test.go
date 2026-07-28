package cfenv_test

import (
	"os"
	"testing"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/go-viper/mapstructure/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testcfenv(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("application deserialization", func() {
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

		when("when not running on Cloud Foundry", func() {
			it("IsRunningOnCF() returns false", func() {
				testEnv := cfenv.Env(notCFEnv)
				_, err := cfenv.New(testEnv)
				Expect(err).To(HaveOccurred())
				Expect(cfenv.IsRunningOnCF()).To(BeFalse())
			})
		})

		when("when running on Cloud Foundry", func() {
			it.Before(func() {
				os.Setenv("VCAP_APPLICATION", "{}")
			})

			it.After(func() {
				os.Unsetenv("VCAP_APPLICATION")
			})

			it("IsRunningOnCF() returns true", func() {
				testEnv := cfenv.Env(cfEnv)
				_, err := cfenv.New(testEnv)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfenv.IsRunningOnCF()).To(BeTrue())
			})
		})

		when("with valid environment", func() {
			it("should deserialize correctly", func() {
				testEnv := cfenv.Env(validEnv)
				env, err := cfenv.New(testEnv)
				Expect(err).To(BeNil())
				Expect(env).NotTo(BeNil())

				Expect(env.ID).To(Equal("451f045fd16427bb99c895a2649b7b2a"))
				Expect(env.InstanceID).To(Equal("451f045fd16427bb99c895a2649b7b2a"))
				Expect(env.AppID).To(Equal("abcabc123123defdef456456"))
				Expect(env.CFAPI).To(Equal("https://api.system_domain.com"))
				Expect(env.Index).To(Equal(0))
				Expect(env.Name).To(Equal("styx-james"))
				Expect(env.SpaceName).To(Equal("jdk"))
				Expect(env.SpaceID).To(Equal("3e0c28c5-6d9c-436b-b9ee-1f4326e54d05"))
				Expect(env.Host).To(Equal("0.0.0.0"))
				Expect(env.Port).To(Equal(61857))
				Expect(env.Version).To(Equal("c1063c1c-40b9-434e-a797-db240b587d32"))
				Expect(env.Home).To(Equal("/home/vcap/app"))
				Expect(env.MemoryLimit).To(Equal("512m"))
				Expect(env.WorkingDir).To(Equal("/home/vcap"))
				Expect(env.TempDir).To(Equal("/home/vcap/tmp"))
				Expect(env.User).To(Equal("vcap"))
				Expect(env.Limits.Disk).To(Equal(1024))
				Expect(env.Limits.Mem).To(Equal(512))
				Expect(env.Limits.FDs).To(Equal(16384))
				Expect(env.ApplicationURIs[0]).To(Equal("styx-james.a1-app.cf-app.com"))
				Expect(len(env.Services)).To(Equal(3))
				Expect(env.Services["elephantsql-dev"][0].Name).To(Equal("elephantsql-dev-c6c60"))
				Expect(env.Services["elephantsql-dev"][0].Label).To(Equal("elephantsql-dev"))
				Expect(env.Services["elephantsql-dev"][0].Tags).To(Equal([]string{"New Product", "relational", "Data Store", "postgresql"}))
				Expect(env.Services["elephantsql-dev"][0].Plan).To(Equal("turtle"))
				Expect(len(env.Services["elephantsql-dev"][0].Credentials)).To(Equal(1))
				Expect(env.Services["elephantsql-dev"][0].Credentials["uri"]).To(Equal("postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"))
				Expect(env.Services["sendgrid"][0].Name).To(Equal("mysendgrid"))
				Expect(env.Services["sendgrid"][0].Label).To(Equal("sendgrid"))
				Expect(env.Services["sendgrid"][0].Tags).To(Equal([]string{"smtp", "Email"}))
				Expect(env.Services["sendgrid"][0].Plan).To(Equal("free"))
				Expect(len(env.Services["sendgrid"][0].Credentials)).To(Equal(3))
				Expect(env.Services["sendgrid"][0].Credentials["hostname"]).To(Equal("smtp.sendgrid.net"))
				Expect(env.Services["sendgrid"][0].Credentials["username"]).To(Equal("QvsXMbJ3rK"))
				Expect(env.Services["sendgrid"][0].Credentials["password"]).To(Equal("HCHMOYluTv"))

				Expect(env.Services["nfs"][0].VolumeMounts[0]["container_dir"]).To(Equal("/testpath"))

				name, err := env.Services.WithName("elephantsql-dev-c6c60")
				Expect(name.Name).To(Equal("elephantsql-dev-c6c60"))
				Expect(err).To(BeNil())

				tag, err := env.Services.WithTag("postgresql")
				Expect(len(tag)).To(Equal(1))
				Expect(tag[0].Tags).To(ContainElement("postgresql"))
				Expect(err).To(BeNil())

				label, err := env.Services.WithLabel("elephantsql-dev")
				Expect(len(label)).To(Equal(1))
				Expect(label[0].Label).To(Equal("elephantsql-dev"))
				Expect(err).To(BeNil())

				names, err := env.Services.WithNameUsingPattern(".*(sql|mysend).*")
				Expect(len(names)).To(Equal(2))
				Expect(err).To(BeNil())
				isValidNames := true
				for _, service := range names {
					if service.Name != "mysendgrid" && service.Name != "elephantsql-dev-c6c60" {
						isValidNames = false
					}
				}
				Expect(isValidNames).To(BeTrue(), "Not valid names when finding by regex")

				tags, err := env.Services.WithTagUsingPattern(".*sql.*")
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

			it("should prefer the PORT environment variable over VCAP_APPLICATION.PORT", func() {
				validEnv = append(validEnv, "PORT=12345")
				testEnv := cfenv.Env(validEnv)
				env, err := cfenv.New(testEnv)
				Expect(err).To(BeNil())
				Expect(env).NotTo(BeNil())
				Expect(env.Port).To(Equal(12345))
			})
		})

		when("without a space name and id", func() {
			it("should deserialize correctly", func() {
				testEnv := cfenv.Env(validEnvWithoutSpaceIDAndName)
				env, err := cfenv.New(testEnv)
				Expect(err).To(BeNil())
				Expect(env).NotTo(BeNil())
				Expect(env.SpaceID).To(BeEmpty())
				Expect(env.SpaceName).To(BeEmpty())
			})
		})

		when("with valid environment with a service with credentials that are an array", func() {
			it("should deserialize correctly", func() {
				testEnv := cfenv.Env(envWithArrayCredentials)
				env, err := cfenv.New(testEnv)
				Expect(err).To(BeNil())
				Expect(env).NotTo(BeNil())

				credential := map[string]interface{}{}
				mapstructure.Decode(env.Services["p-kafka"][0].Credentials["kafka"], &credential)

				Expect(len(env.Services["p-kafka"][0].Credentials)).To(Equal(1))
				Expect(credential["node_ips"]).To(Equal([]interface{}{"10.244.9.2", "10.244.9.6", "10.244.9.10"}))
				Expect(credential["port"]).To(Equal(float64(9092)))
			})
		})

		when("with valid environment with a service with credentials with a port that is an int", func() {
			it("should to deserialize correctly", func() {
				testEnv := cfenv.Env(envWithIntCredentials)
				env, err := cfenv.New(testEnv)
				Expect(err).To(BeNil())
				Expect(env).NotTo(BeNil())

				Expect(env.ID).To(Equal("451f045fd16427bb99c895a2649b7b2a"))
				Expect(env.Index).To(Equal(0))
				Expect(env.Name).To(Equal("styx-james"))
				Expect(env.Host).To(Equal("0.0.0.0"))
				Expect(env.Port).To(Equal(61857))
				Expect(env.Version).To(Equal("c1063c1c-40b9-434e-a797-db240b587d32"))
				Expect(env.Home).To(Equal("/home/vcap/app"))
				Expect(env.MemoryLimit).To(Equal("512m"))
				Expect(env.WorkingDir).To(Equal("/home/vcap"))
				Expect(env.TempDir).To(Equal("/home/vcap/tmp"))
				Expect(env.User).To(Equal("vcap"))
				Expect(env.ApplicationURIs[0]).To(Equal("styx-james.a1-app.cf-app.com"))
				Expect(len(env.Services)).To(Equal(3))

				Expect(env.Services["elephantsql-dev"][0].Name).To(Equal("elephantsql-dev-c6c60"))
				Expect(env.Services["elephantsql-dev"][0].Label).To(Equal("elephantsql-dev"))
				Expect(env.Services["elephantsql-dev"][0].Tags).To(Equal([]string{"New Product", "relational", "Data Store", "postgresql"}))
				Expect(env.Services["elephantsql-dev"][0].Plan).To(Equal("turtle"))
				Expect(len(env.Services["elephantsql-dev"][0].Credentials)).To(Equal(1))
				Expect(env.Services["elephantsql-dev"][0].Credentials["uri"]).To(Equal("postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"))

				Expect(env.Services["cloudantNoSQLDB"][0].Name).To(Equal("my_cloudant"))
				Expect(env.Services["cloudantNoSQLDB"][0].Label).To(Equal("cloudantNoSQLDB"))
				Expect(env.Services["cloudantNoSQLDB"][0].Plan).To(Equal("Shared"))
				Expect(len(env.Services["cloudantNoSQLDB"][0].Credentials)).To(Equal(5))
				Expect(env.Services["cloudantNoSQLDB"][0].Credentials["port"]).To(Equal(float64(443)))

				Expect(env.Services["sendgrid"][0].Name).To(Equal("mysendgrid"))
				Expect(env.Services["sendgrid"][0].Label).To(Equal("sendgrid"))
				Expect(env.Services["sendgrid"][0].Tags).To(Equal([]string{"smtp", "Email"}))
				Expect(env.Services["sendgrid"][0].Plan).To(Equal("free"))
				Expect(len(env.Services["sendgrid"][0].Credentials)).To(Equal(3))
				Expect(env.Services["sendgrid"][0].Credentials["hostname"]).To(Equal("smtp.sendgrid.net"))
				Expect(env.Services["sendgrid"][0].Credentials["username"]).To(Equal("QvsXMbJ3rK"))
				Expect(env.Services["sendgrid"][0].Credentials["password"]).To(Equal("HCHMOYluTv"))

				name, err := env.Services.WithName("elephantsql-dev-c6c60")
				Expect(name.Name).To(Equal("elephantsql-dev-c6c60"))
				Expect(err).To(BeNil())

				tag, err := env.Services.WithTag("postgresql")
				Expect(len(tag)).To(Equal(1))
				Expect(tag[0].Tags).To(ContainElement("postgresql"))
				Expect(err).To(BeNil())

				label, err := env.Services.WithLabel("elephantsql-dev")
				Expect(len(label)).To(Equal(1))
				Expect(label[0].Label).To(Equal("elephantsql-dev"))
				Expect(err).To(BeNil())

				names, err := env.Services.WithNameUsingPattern(".*(sql|my_cloud).*")
				Expect(len(names)).To(Equal(2))
				Expect(err).To(BeNil())
				isValidNames := true
				for _, service := range names {
					if service.Name != "my_cloudant" && service.Name != "elephantsql-dev-c6c60" {
						isValidNames = false
					}
				}
				Expect(isValidNames).To(BeTrue(), "Not valid names when finding by regex")

				tags, err := env.Services.WithTagUsingPattern(".*s.*")
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

		when("with invalid environment", func() {
			it("should deserialize correctly, with missing values", func() {
				testEnv := cfenv.Env(invalidEnv)
				env, err := cfenv.New(testEnv)
				Expect(err).To(BeNil())
				Expect(env).NotTo(BeNil())

				Expect(env.ID).To(Equal(""))
				Expect(env.Index).To(Equal(0))
				Expect(env.Name).To(Equal("styx-james"))
				Expect(env.Host).To(Equal("0.0.0.0"))
				Expect(env.Port).To(Equal(61857))
				Expect(env.Version).To(Equal("c1063c1c-40b9-434e-a797-db240b587d32"))
				Expect(env.Home).To(Equal("/home/vcap/app"))
				Expect(env.MemoryLimit).To(Equal(""))
				Expect(env.WorkingDir).To(Equal("/home/vcap"))
				Expect(env.TempDir).To(Equal("/home/vcap/tmp"))
				Expect(env.User).To(Equal("vcap"))
				Expect(env.ApplicationURIs[0]).To(Equal("styx-james.a1-app.cf-app.com"))
				Expect(len(env.Services)).To(Equal(2))
				Expect(len(env.Services)).To(Equal(2))
				Expect(env.Services["elephantsql-dev"][0].Name).To(Equal(""))
				Expect(env.Services["elephantsql-dev"][0].Label).To(Equal("elephantsql-dev"))
				Expect(env.Services["elephantsql-dev"][0].Plan).To(Equal("turtle"))
				Expect(len(env.Services["elephantsql-dev"][0].Credentials)).To(Equal(1))
				Expect(env.Services["elephantsql-dev"][0].Credentials["uri"]).To(Equal("postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@babar.elephantsql.com:5432/seilbmbd"))

				Expect(env.Services["sendgrid"][0].Name).To(Equal("mysendgrid"))
				Expect(env.Services["sendgrid"][0].Label).To(Equal("sendgrid"))
				Expect(env.Services["sendgrid"][0].Plan).To(Equal("free"))
				Expect(len(env.Services["sendgrid"][0].Credentials)).To(Equal(3))
				Expect(env.Services["sendgrid"][0].Credentials["hostname"]).To(Equal("smtp.sendgrid.net"))
				Expect(env.Services["sendgrid"][0].Credentials["username"]).To(Equal("QvsXMbJ3rK"))
				Expect(env.Services["sendgrid"][0].Credentials["password"]).To(Equal("HCHMOYluTv"))
			})
		})
	})

	when("CredentialString", func() {
		var service = cfenv.Service{
			Credentials: map[string]interface{}{
				"string": "stringy-credential",
				"int":    42,
				"nested": map[string]string{
					"key": "value",
				},
			},
		}

		it("returns the requested credential as a string when the credential is a string", func() {
			result, ok := service.CredentialString("string")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal("stringy-credential"))
		})

		it("returns false when the credential is not a string", func() {
			_, ok := service.CredentialString("int")
			Expect(ok).To(BeFalse())
		})

		it("returns false when the credential is a nested thing", func() {
			_, ok := service.CredentialString("nested")
			Expect(ok).To(BeFalse())
		})
	})

	// Issue #11: a credential whose value was a nested object once panicked
	// the decoder, which expected every credential to be a string. Pins the
	// payload from that report so the fix cannot silently regress.
	when("credentials containing nested objects", func() {
		it("decodes them and leaves the nested object reachable", func() {
			env, err := cfenv.New(map[string]string{
				"VCAP_APPLICATION": "{}",
				"VCAP_SERVICES": `{"redis-lite":[{
					"credentials":{
						"hostname":"10.10.10.1",
						"password":"abc",
						"port":"32843",
						"ports":{"6379/tcp":"32843"}
					},
					"label":"redis-lite",
					"name":"uptime-redis",
					"plan":"free",
					"tags":["redis28","redis","key-value"]
				}]}`,
			})
			Expect(err).To(BeNil())

			service, err := env.Services.WithName("uptime-redis")
			Expect(err).To(BeNil())

			port, ok := service.Credential("ports", "6379/tcp")
			Expect(ok).To(BeTrue())
			Expect(port).To(Equal("32843"))
		})
	})

	when("Credential", func() {
		// Mirrors what cfenv.New produces: nested objects decode to
		// map[string]interface{} and leaves keep their JSON types.
		var service = cfenv.Service{
			Credentials: map[string]interface{}{
				"uri": "amqp://guest@rabbitmq:5672",
				"protocols": map[string]interface{}{
					"amqp": map[string]interface{}{
						"uri":  "amqp://guest:guest@rabbitmq:5672",
						"ssl":  false,
						"port": float64(5672),
					},
				},
				"hosts": []interface{}{"a", "b"},
				// A real credential key that contains a dot.
				"jdbc.url": "jdbc:mysql://x",
			},
		}

		it("returns a top-level credential", func() {
			result, ok := service.Credential("uri")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal("amqp://guest@rabbitmq:5672"))
		})

		it("returns a credential nested several levels deep", func() {
			result, ok := service.Credential("protocols", "amqp", "uri")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal("amqp://guest:guest@rabbitmq:5672"))
		})

		it("returns a non-string leaf with its own type", func() {
			result, ok := service.Credential("protocols", "amqp", "ssl")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal(false))

			result, ok = service.Credential("protocols", "amqp", "port")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal(float64(5672)))
		})

		it("returns an intermediate map when the path stops there", func() {
			result, ok := service.Credential("protocols", "amqp")
			Expect(ok).To(BeTrue())
			Expect(result).To(HaveKeyWithValue("uri", "amqp://guest:guest@rabbitmq:5672"))
		})

		it("returns a slice leaf", func() {
			result, ok := service.Credential("hosts")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal([]interface{}{"a", "b"}))
		})

		it("addresses a key containing a dot", func() {
			result, ok := service.Credential("jdbc.url")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal("jdbc:mysql://x"))
		})

		it("returns false when a top-level key is missing", func() {
			_, ok := service.Credential("nope")
			Expect(ok).To(BeFalse())
		})

		it("returns false when a nested key is missing", func() {
			_, ok := service.Credential("protocols", "amqp", "nope")
			Expect(ok).To(BeFalse())
		})

		it("returns false when descending through a non-map", func() {
			_, ok := service.Credential("uri", "nope")
			Expect(ok).To(BeFalse())
		})

		it("returns false when given no keys", func() {
			_, ok := service.Credential()
			Expect(ok).To(BeFalse())
		})

		it("returns false when the credentials are empty", func() {
			empty := cfenv.Service{}
			_, ok := empty.Credential("uri")
			Expect(ok).To(BeFalse())
		})
	})

	when("CredentialPath", func() {
		var service = cfenv.Service{
			Credentials: map[string]interface{}{
				"protocols": map[string]interface{}{
					"amqp": map[string]interface{}{
						"uri": "amqp://guest:guest@rabbitmq:5672",
						"ssl": false,
					},
				},
				"jdbc.url": "jdbc:mysql://x",
			},
		}

		it("returns a credential addressed by a dotted path", func() {
			result, ok := service.CredentialPath("protocols.amqp.uri")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal("amqp://guest:guest@rabbitmq:5672"))
		})

		it("returns a non-string leaf with its own type", func() {
			result, ok := service.CredentialPath("protocols.amqp.ssl")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal(false))
		})

		it("returns false for a missing path", func() {
			_, ok := service.CredentialPath("protocols.stomp.uri")
			Expect(ok).To(BeFalse())
		})

		it("returns false for an empty path", func() {
			_, ok := service.CredentialPath("")
			Expect(ok).To(BeFalse())
		})

		// The documented limitation of the dotted form: a key containing a
		// dot is split and becomes unreachable. Credential addresses it.
		it("cannot address a key containing a dot", func() {
			_, ok := service.CredentialPath("jdbc.url")
			Expect(ok).To(BeFalse())

			result, ok := service.Credential("jdbc.url")
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal("jdbc:mysql://x"))
		})
	})
}
