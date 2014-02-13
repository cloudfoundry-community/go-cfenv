package cfenv

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

func New(env map[string]string) *App {
	var app App
	appVar := env["VCAP_APPLICATION"]
	if err := json.Unmarshal([]byte(appVar), &app); err != nil {
		// panic(err)
	}
	app.Home = env["HOME"]
	app.MemoryLimit = env["MEMORY_LIMIT"]
	app.WorkingDir = env["PWD"]
	app.TempDir = env["TMPDIR"]
	app.User = env["USER"]
	var rawServices map[string]interface{}
	servicesVar := env["VCAP_SERVICES"]
	if err := json.Unmarshal([]byte(servicesVar), &rawServices); err != nil {
		// panic(err)
	}

	services := make(map[string][]Service)
	for k, v := range rawServices {
		fmt.Println(v)
		var serviceInstances []Service
		if err := mapstructure.Decode(v, &serviceInstances); err != nil {
			// panic(err)
		}
		fmt.Println(serviceInstances[0].Name)
		fmt.Println(serviceInstances[0].Label)
		fmt.Println(serviceInstances[0].Plan)
		services[k] = serviceInstances
	}
	app.Services = services
	return &app
}

func Current() *App {
	return New(CurrentEnv())
}
