package cfenv

// For bindable services Cloud Foundry will add connection details to the VCAP_SERVICES environment variable when you restart your application, after binding a service instance to your application.
// The results are returned as a JSON document that contains an object for each service for which one or more instances are bound to the application. The service object contains a child object for each service instance of that service that is bound to the application. The attributes that describe a bound service are represented in the Service struct.
type Service struct {
	Name        string            // name of the service
	Label       string            // label of the service
	Plan        string            // plan of the service
	Credentials map[string]string // credentials for the service
}
