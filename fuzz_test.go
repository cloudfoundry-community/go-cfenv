package cfenv_test

import (
	"testing"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

// FuzzNew exercises the two documents New parses. Their contents originate
// with the platform and the service brokers bound to an app, not with the
// app itself, so this package is the boundary that has to survive whatever
// arrives. New is expected to return an app or an error for any input, and
// never to panic.
func FuzzNew(f *testing.F) {
	f.Add(`{}`, `{}`)
	f.Add(
		`{"instance_id":"abc","name":"app","limits":{"mem":512,"disk":1024,"fds":16384}}`,
		`{"redis":[{"name":"r","label":"redis","tags":["t"],"plan":"free","credentials":{"port":"6379","ports":{"6379/tcp":"32843"}}}]}`,
	)
	f.Add(`{}`, `{"user-provided":[{"name":"ups","syslog_drain_url":"syslog://x:514","credentials":{"a":{"b":{"c":1}}}}]}`)
	f.Add(`{}`, `{"nfs":[{"volume_mounts":[{"container_dir":"/d","mode":"rw"}]}]}`)
	f.Add(`{"limits":null}`, `{"x":[]}`)
	f.Add(``, ``)

	f.Fuzz(func(t *testing.T, application, services string) {
		app, err := cfenv.New(map[string]string{
			"VCAP_APPLICATION": application,
			"VCAP_SERVICES":    services,
		})
		if err != nil {
			if app != nil {
				t.Fatalf("returned both an app and an error: %v", err)
			}

			return
		}

		if app == nil {
			t.Fatal("returned neither an app nor an error")
		}

		// Whatever decoded, the accessors have to stay total: they may find
		// nothing, but they must not panic on any shape that got this far.
		for _, instances := range app.Services {
			for i := range instances {
				service := instances[i]
				service.CredentialString("uri")
				service.Credential("a", "b", "c")
				service.CredentialPath("a.b.c")
			}
		}

		app.Services.WithTag("t")
		app.Services.WithLabel("redis")
		app.Services.WithName("r")
		app.Services.WithNameUsingPattern("r.*")
		app.Services.WithTagUsingPattern("t.*")
	})
}
