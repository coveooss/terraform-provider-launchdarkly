package launchdarkly

import (
	"regexp"
	"errors"
	"fmt"
)

func validateKey(v interface{}, k string) ([]string, []error) {
	value := v.(string)

	if len(value) < 1 || len(value) > 20 {
		return nil, []error{errors.New(fmt.Sprintf("%s must be between 1 and 20 characters: %s", k, value))}
	}

	matched, err := regexp.MatchString("^[A-Za-z0-9_\\-\\.]+$", value)
	if err != nil {
		return nil, []error{err}
	}

	if !matched {
		return nil, []error{errors.New(fmt.Sprintf("%s is not a valid key: %s", k, value))}
	}

	return nil, nil
}

func validateColor(v interface{}, k string) ([]string, []error) {
	value := v.(string)

	matched, err := regexp.MatchString("^[0-9a-fA-F]{6}$", value)
	if err != nil {
		return nil, []error{err}
	}

	if !matched {
		return nil, []error{errors.New(fmt.Sprintf("%s is not a valid RGB color code: %s", k, value))}
	}

	return nil, nil
}
