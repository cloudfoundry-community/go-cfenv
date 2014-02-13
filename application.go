package cfenv

type App struct {
	Id      string `json:"instance_id"`    // id of the app
	Index   int    `json:"instance_index"` // index of the app
	Name    string `json:"name"`           // name of the app
	Host    string `json:"host"`           // host of the app
	Port    int    `json:"port"`           // port of the app
	Version string `json:"version"`        // version of the app
}
