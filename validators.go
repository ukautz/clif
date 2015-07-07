package clif

import (
	"fmt"
	"regexp"
)

// IsAny joins a set of validator methods and returns true if ANY of them matches
func IsAny(v ...ParseMethod) ParseMethod {
	return func(name, value string) (string, error) {
		var err error
		replace := value
		for _, c := range v {
			if replace, err = c(name, replace); err == nil {
				return replace, nil
			}
		}
		return "", err
	}
}

// IsAll joins a set of validators methods and returns true if ALL of them match
func IsAll(v ...ParseMethod) ParseMethod {
	return func(name, value string) (string, error) {
		var err error
		replace := value
		for _, c := range v {
			if replace, err = c(name, replace); err != nil {
				return "", err
			}
		}
		return replace, nil
	}
}

var rxIsInt = regexp.MustCompile(`^[1-9][0-9]*$`)

// IsInt checks if value is an integer
func IsInt(name, value string) (string, error) {
	if !rxIsInt.MatchString(value) {
		return "", fmt.Errorf("Is not integer")
	}
	return value, nil
}

var rxIsFloat = regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)?$`)

// IsFloat checks if value is float
func IsFloat(name, value string) (string, error) {
	if !rxIsFloat.MatchString(value) {
		return "", fmt.Errorf("Is not float")
	}
	return value, nil
}
