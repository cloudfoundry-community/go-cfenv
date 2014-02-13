package cfenv

type Service struct {
	Name        string            //`mapstructure:"squash,name"`  // name of the service
	Label       string            //`mapstructure:"squash,label"` // label of the service
	Plan        string            //`mapstructure:"squash,plan"`  // plan of the service
	Credentials map[string]string //`json:"credentials"`          // credentials for the service
}
