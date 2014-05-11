package cfenv

import (
	"fmt"
	"strings"
)

// Service describes a bound service. For bindable services Cloud Foundry will
// add connection details to the VCAP_SERVICES environment variable when you
// restart your application, after binding a service instance to your
// application.
//
// The results are returned as a JSON document that contains an object for each
// service for which one or more instances are bound to the application. The
// service object contains a child object for each service instance of that
// service that is bound to the application.
type Service struct {
	Name        string            // name of the service
	Label       string            // label of the service
	Tags        []string          // tags for the service
	Plan        string            // plan of the service
	Credentials map[string]string // credentials for the service
}

// Services is a map associating service labels and an array of services with
// that label.
type Services map[string][]Service

// WithTag find services with the specified tag
func (s *Services) WithTag(tag string) ([]Service, error) {
	r := []Service{}
	for _, v := range *s {
		for i := range v {
			service := v[i]
			for t := range service.Tags {
				if strings.EqualFold(tag, service.Tags[t]) {
					r = append(r, service)
				}
			}
		}
	}

	if len(r) > 0 {
		return r, nil
	}

	return nil, fmt.Errorf("no services with tag %s", tag)
}

// WithLabel finds the service with the specified label
func (s *Services) WithLabel(label string) ([]Service, error) {
	for k, v := range *s {
		if strings.EqualFold(label, k) {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no services with label %s", label)
}

// WithName finds the service with the specified name
func (s *Services) WithName(name string) (*Service, error) {
	for _, v := range *s {
		for i := range v {
			service := v[i]
			if strings.EqualFold(name, service.Name) {
				return &service, nil
			}
		}
	}

	return nil, fmt.Errorf("no service with name %s", name)
}
