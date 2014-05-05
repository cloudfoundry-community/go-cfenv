package cfenv

import (
	"errors"
)

// Service describes a bound service. For bindable services Cloud Foundry will add connection details to the VCAP_SERVICES environment variable when you restart your application, after binding a service instance to your application.
// The results are returned as a JSON document that contains an object for each service for which one or more instances are bound to the application. The service object contains a child object for each service instance of that service that is bound to the application.
type Service struct {
	Name        string            // name of the service
	Label       string            // label of the service
	Tags        []string          // tags for the service
	Plan        string            // plan of the service
	Credentials map[string]string // credentials for the service
}

type Services map[string][]Service

func (s Services) FindByTagName(tag string) (Service, error) {
	for _, services := range s {
		for serviceIndex := range services {
			service := services[serviceIndex]
			for tagIndex := range service.Tags {
				if tag == service.Tags[tagIndex] {
					return service, nil
				}
			}
		}
	}
	return Service{}, errors.New("Error finding service by tag")
}

func (s Services) FindByLabel(label string) (Service, error) {
	for _, services := range s {
		for serviceIndex := range services {
			service := services[serviceIndex]
			if label == service.Label {
				return service, nil
			}
		}
	}
	return Service{}, errors.New("Error finding service by label")
}

func (s Services) FindByName(name string) (Service, error) {
	for _, services := range s {
		for serviceIndex := range services {
			service := services[serviceIndex]
			if name == service.Name {
				return service, nil
			}
		}
	}
	return Service{}, errors.New("Error finding service by name")
}
