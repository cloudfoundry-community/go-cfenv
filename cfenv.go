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
	app.Home = env["HOME"]
	app.MemoryLimit = env["MEMORY_LIMIT"]
	app.WorkingDir = env["PWD"]
	app.TempDir = env["TMPDIR"]
	app.User = env["USER"]
	return &app
}

func Current() *App {
	return New(CurrentEnv())
}
