package cfenv

import (
	"fmt"
	"regexp"
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
	Name         string                 // the binding name if the binding has one, otherwise the instance name
	Label        string                 // name of the service offering; "user-provided" for a user-provided instance
	Tags         []string               // tags for the service
	Plan         string                 // plan of the service
	Credentials  map[string]interface{} // credentials for the service
	VolumeMounts []map[string]string    `mapstructure:"volume_mounts"` // volume mount info as provided by the nfsbroker

	SyslogDrainURL string `mapstructure:"syslog_drain_url"` // where the service wants app logs streamed; set by `cf cups -l`
	InstanceGUID   string `mapstructure:"instance_guid"`    // GUID of the service instance
	InstanceName   string `mapstructure:"instance_name"`    // name the user gave the service instance
	BindingGUID    string `mapstructure:"binding_guid"`     // GUID of the service binding
	BindingName    string `mapstructure:"binding_name"`     // name the user gave the binding, if any
}

func (s *Service) CredentialString(key string) (string, bool) {
	credential, ok := s.Credentials[key].(string)

	return credential, ok
}

// Credential walks the credentials one key at a time and returns the value at
// the end of the path, whatever its type: a string, a bool, the float64 a JSON
// number decodes to, or a nested map or slice.
//
// Passing the keys separately means every key is addressable, including one
// that contains a dot:
//
//	uri, ok := service.Credential("protocols", "amqp", "uri")
//	url, ok := service.Credential("jdbc.url")
//
// It reports false if any key along the path is absent, if the path descends
// through a value that is not a nested object, or if no keys are given.
func (s *Service) Credential(keys ...string) (interface{}, bool) {
	if len(keys) == 0 {
		return nil, false
	}

	var current interface{} = s.Credentials

	for _, key := range keys {
		object, isObject := current.(map[string]interface{})
		if !isObject {
			return nil, false
		}

		value, exists := object[key]
		if !exists {
			return nil, false
		}

		current = value
	}

	return current, true
}

// CredentialPath is Credential addressed by a single dot-delimited path:
//
//	uri, ok := service.CredentialPath("protocols.amqp.uri")
//
// Because the path is split on every dot, it cannot address a key that itself
// contains one — CredentialPath("jdbc.url") looks for "url" inside "jdbc" and
// reports false. Use Credential for those keys.
func (s *Service) CredentialPath(path string) (interface{}, bool) {
	if path == "" {
		return nil, false
	}

	return s.Credential(strings.Split(path, ".")...)
}

// Services is an association of service labels to a slice of services with that
// label.
type Services map[string][]Service

// WithTag finds services with the specified tag.
func (s *Services) WithTag(tag string) ([]Service, error) {
	result := []Service{}

	for _, services := range *s {
		for i := range services {
			service := services[i]
			for _, t := range service.Tags {
				if strings.EqualFold(tag, t) {
					result = append(result, service)

					break
				}
			}
		}
	}

	if len(result) > 0 {
		return result, nil
	}

	return nil, fmt.Errorf("no services with tag %s", tag)
}

// WithTag finds services with a tag pattern.
func (s *Services) WithTagUsingPattern(tagPattern string) ([]Service, error) {
	result := []Service{}

	for _, services := range *s {
		for i := range services {
			service := services[i]
			for _, t := range service.Tags {
				if s.match(tagPattern, t) {
					result = append(result, service)

					break
				}
			}
		}
	}

	if len(result) > 0 {
		return result, nil
	}

	return nil, fmt.Errorf("no services with tag pattern %s", tagPattern)
}

// WithLabel finds the service with the specified label.
func (s *Services) WithLabel(label string) ([]Service, error) {
	for l, services := range *s {
		if strings.EqualFold(label, l) {
			return services, nil
		}
	}

	return nil, fmt.Errorf("no services with label %s", label)
}

// WithName finds the service with a name pattern.
func (s *Services) WithNameUsingPattern(namePattern string) ([]Service, error) {
	result := []Service{}

	for _, services := range *s {
		for i := range services {
			service := services[i]
			if s.match(namePattern, service.Name) {
				result = append(result, service)
			}
		}
	}

	if len(result) > 0 {
		return result, nil
	}

	return nil, fmt.Errorf("no service with name pattern %s", namePattern)
}

// WithName finds the service with the specified name.
func (s *Services) WithName(name string) (*Service, error) {
	for _, services := range *s {
		for i := range services {
			service := services[i]
			if strings.EqualFold(name, service.Name) {
				return &service, nil
			}
		}
	}

	return nil, fmt.Errorf("no service with name %s", name)
}

func (s *Services) match(matcher, content string) bool {
	regex, err := regexp.Compile("(?i)^" + matcher + "$")
	if err != nil {
		return false
	}

	return regex.MatchString(content)
}
