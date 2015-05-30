package cli

import (
	"fmt"
	"regexp"
)

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

func IsInt(name, value string) error {
	if !rxIsInt.MatchString(value) {
		return fmt.Errorf("Is not integer")
	}
	return nil
}

var rxIsFloat = regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)?$`)

func IsFloat(name, value string) error {
	if !rxIsFloat.MatchString(value) {
		return fmt.Errorf("Is not float")
	}
	return nil
}
