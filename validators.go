package clif

import (
	"fmt"
	"regexp"
)

// IsAny joins a set of validator methods and returns true if ANY of them matches
func IsAny(v ...ValidatorMethod) ValidatorMethod {
	return func(name, value string) error {
		var err error
		for _, c := range v {
			if err = c(name, value); err == nil {
				return nil
			}
		}
		return err
	}
}

// IsAll joins a set of validators methods and returns true if ALL of them match
func IsAll(v ...ValidatorMethod) ValidatorMethod {
	return func(name, value string) error {
		for _, c := range v {
			if err := c(name, value); err != nil {
				return err
			}
		}
		return nil
	}
}

var rxIsInt = regexp.MustCompile(`^[1-9][0-9]*$`)

// IsInt checks if value is an integer
func IsInt(name, value string) error {
	if !rxIsInt.MatchString(value) {
		return fmt.Errorf("Is not integer")
	}
	return nil
}

var rxIsFloat = regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)?$`)

// IsInt checks if value is float
func IsFloat(name, value string) error {
	if !rxIsFloat.MatchString(value) {
		return fmt.Errorf("Is not float")
	}
	return nil
}
