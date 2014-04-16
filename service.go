package cfenv

// Service describes a bound service. For bindable services Cloud Foundry will add connection details to the VCAP_SERVICES environment variable when you restart your application, after binding a service instance to your application.
// The results are returned as a JSON document that contains an object for each service for which one or more instances are bound to the application. The service object contains a child object for each service instance of that service that is bound to the application.
type Service struct {
	Name        string            // name of the service
	Label       string            // label of the service
	Tags        []string          // tags for the service
	Plan        string            // plan of the service
	Credentials map[string]string // credentials for the service
}
