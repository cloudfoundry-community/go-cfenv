package cfenv

import (
	"encoding/json"
)

func New(env map[string]string) *App {
	var app App
	appVar := env["VCAP_APPLICATION"]
	if err := json.Unmarshal([]byte(appVar), &app); err != nil {
		// panic(err)
	}
	return &app
}

func Current() *App {
	return New(CurrentEnv())
}
