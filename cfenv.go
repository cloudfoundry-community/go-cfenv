// Package cfenv provides information about the current app deployed on Cloud Foundry, including any bound service(s).
package cfenv

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

// New creates a new App with the provided environment.
func New(env map[string]string) (*App, error) {
	var app App

	appVar := env["VCAP_APPLICATION"]
	if err := json.Unmarshal([]byte(appVar), &app); err != nil {
		return nil, err
	}
	// duplicate the InstanceID to the previously named ID field for backwards
	// compatibility
	app.ID = app.InstanceID

	app.Home = env["HOME"]
	app.MemoryLimit = env["MEMORY_LIMIT"]

	if port, err := strconv.Atoi(env["PORT"]); err == nil {
		app.Port = port
	}

	app.WorkingDir = env["PWD"]
	app.TempDir = env["TMPDIR"]
	app.User = env["USER"]

	servicesDoc, err := servicesJSON(env)
	if err != nil {
		return nil, err
	}

	var rawServices map[string]interface{}

	if err := json.Unmarshal(servicesDoc, &rawServices); err != nil {
		return nil, err
	}

	services := make(map[string][]Service)

	for key, value := range rawServices {
		var serviceInstances []Service
		if err := mapstructure.WeakDecode(value, &serviceInstances); err != nil {
			return nil, err
		}

		services[key] = serviceInstances
	}

	app.Services = services

	return &app, nil
}

// servicesJSON returns the raw service bindings document, read from the file
// named by VCAP_SERVICES_FILE_PATH when that is set and from VCAP_SERVICES
// otherwise.
//
// Cloud Foundry sets one or the other, never both: with the
// file-based-vcap-services app feature enabled it provides only the file path,
// which is how bindings too large for an environment variable are delivered
// (RFC-0030). The file wins here, so the tie only matters in a synthetic
// environment.
//
// A file path that cannot be read is an error rather than a fall back to
// VCAP_SERVICES. On a file-based app that variable is unset, so falling back
// would hand the caller an app with no services and no explanation.
func servicesJSON(env map[string]string) ([]byte, error) {
	path := env["VCAP_SERVICES_FILE_PATH"]
	if path == "" {
		return []byte(env["VCAP_SERVICES"]), nil
	}

	// The whole point of the feature is to read the path the platform hands
	// us, so the path is necessarily a variable.
	// #nosec G304
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading VCAP_SERVICES_FILE_PATH %q: %w", path, err)
	}

	return contents, nil
}

// Current creates a new App with the current environment; returns an error if the current environment is not a Cloud Foundry environment.
func Current() (*App, error) {
	return New(CurrentEnv())
}

// IsRunningOnCF returns true if the current environment is Cloud Foundry and false if it is not Cloud Foundry.
func IsRunningOnCF() bool {
	return strings.TrimSpace(os.Getenv("VCAP_APPLICATION")) != ""
}
