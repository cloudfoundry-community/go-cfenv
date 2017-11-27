package cfenv

import (
	"fmt"
	"regexp"
	"strconv"
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
	Name        string                 // name of the service
	Label       string                 // label of the service
	Tags        []string               // tags for the service
	Plan        string                 // plan of the service
	Credentials map[string]interface{} `json:"credentials"` // credentials for the service
}

func (s *Service) CredentialString(key string) (string, bool) {
	credential, ok := s.Credentials[key].(string)
	return credential, ok
}

// CredentialInterface get specified credentail in a service
func (s *Service) Credential(key string) (interface{}, bool) {
	var o interface{}

	o = s.Credentials
	for _, p := range strings.Split(key, ".") {
		switch o.(type) {
		case map[string]interface{}:
			v, ok := o.(map[string]interface{})[p]
			if !ok {
				return nil, false
			}
			o = valueType(v)

		case []interface{}:
			u, err := strconv.ParseUint(p, 10, 0)
			if err != nil {
				return nil, false
			}
			i := int(u)
			if i >= len(o.([]interface{})) {
				return nil, false
			}
			o = valueType(o.([]interface{})[i])

		default:
			return nil, false
		}
	}
	return o, true
}

func valueType(v interface{}) (o interface{}) {
	o = v
	switch v.(type) {
	case float64:
		if v == float64(int64(v.(float64))) {
			o = int64(v.(float64))
		}

	case string:
		l := strings.ToLower(v.(string))
		switch l {
		case "true":
			o = true
		case "false":
			o = false
		}
	}
	return
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
func (s *Services) match(matcher, content string) bool {
	regex, err := regexp.Compile("(?i)^" + matcher + "$")
	if err != nil {
		return false
	}
	return regex.MatchString(content)
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
