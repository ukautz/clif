package clif

import (
	"encoding/json"
	"fmt"
	"github.com/ukautz/reflekt"
	"regexp"
	//	"strings"
	"time"
)

/*
Options and Argument are parameters for commands.

Arguments are fixed positioned, meaning their order does matter. Options

	command foo --bar baz -ding

*/

// SetupMethod is type for callback on Setup of Argument or Option. The return
// string replaces the original input. Called on each value (in case of multiple).
type ParseMethod func(name, value string) (string, error)

// parameter is core for Argument and Option
type parameter struct {

	// Name is used for describing and accessing this argument
	Name string

	// Value holds what was provided on the command line
	Values []string

	// Usage is a short description of this argument
	Usage string

	// Descriptions is a lengthy elaboration of the purpose, use-case, life-story of this argument
	Description string

	// Required determines whether command can execute without this argument
	// Should NOT be changed after adding with `AddCommand` from `Command`
	Required bool

	// Multiple decides whether multiple values are allowed.
	Multiple bool

	// Default is used if no value is provided
	Default string

	// Optional environment variable name, which is used in case not provided (overwrites
	// any default)
	Env string

	// Parse is optional callback, which is applied on parameter values
	// they are assigned. It can be used to validate, transform or otherwise
	// utilize user provided inputs. Mind that inputs can be multiple and it
	// will be called for each of those multiple inputs.
	Parse ParseMethod

	// Regex for checking if input value can be accepted
	Regex *regexp.Regexp
}

/*
---------------------
SETTER
---------------------
*/

// Assign tries to add value to parameter and returns error if it fails due to invalid format or
// invalid amount (single vs multiple parameters)
func (this *parameter) Assign(val string) error {
	if this.Values == nil {
		this.Values = make([]string, 0)
	}
	l := len(this.Values)
	if l > 0 && !this.Multiple {
		return fmt.Errorf("Parameter \"%s\" does not support multiple values", this.Name)
	} else {
		print := func(m string) string {
			return fmt.Sprintf("Parameter \"%s\" invalid: %s", this.Name, m)
		}
		if l > 1 {
			print = func(m string) string {
				return fmt.Sprintf("Parameter \"%s\" (%d) is invalid: %s", this.Name, l+2, m)
			}
		}
		if this.Regex != nil && !this.Regex.MatchString(val) {
			return fmt.Errorf(print("Does not match criteria"))
		}
		if this.Parse != nil {
			if replace, err := this.Parse(this.Name, val); err != nil {
				return fmt.Errorf(print(err.Error()))
			} else {
				val = replace
			}
		}
		this.Values = append(this.Values, val)
		return nil
	}
}

/*
---------------------
GETTER
---------------------
*/

// Provided returns bool whether argument was provided
func (this *parameter) Provided() bool {
	return this.Values != nil
}

// Provided returns amount of values provided
func (this *parameter) Count() int {
	return len(this.Values)
}

// String representation of the value (can be empty string)
func (this *parameter) String() string {
	if this.Values == nil {
		return ""
	} else {
		return this.Values[0]
	}
}

// Strings returns values as array of strings
func (this *parameter) Strings() []string {
	return this.Values
}

// Int representation of the value (will be 0, if not given or not parsable)
func (this *parameter) Int() int {
	if this.Values == nil {
		return 0
	} else {
		return reflekt.AsInt(this.Values[0])
	}
}

// Ints returns values as int array (values will be 0, if not parsable to int)
func (this *parameter) Ints() []int {
	if this.Values == nil {
		return nil
	} else {
		res := make([]int, this.Count())
		for i, v := range this.Values {
			res[i] = reflekt.AsInt(v)
		}
		return res
	}
}

// Float representation of the value (will be 0.0, if not given or not parsable)
func (this *parameter) Float() float64 {
	if this.Values == nil {
		return 0
	} else {
		return reflekt.AsFloat(this.Values[0])
	}
}

// Floats returns values as float64 array (values will be 0.0, if not parsable to float64)
func (this *parameter) Floats() []float64 {
	if this.Values == nil {
		return nil
	} else {
		res := make([]float64, this.Count())
		for i, v := range this.Values {
			res[i] = reflekt.AsFloat(v)
		}
		return res
	}
}

// Bool representation of the value (will be false, if not given or not parsable)
func (this *parameter) Bool() bool {
	if this.Values == nil {
		return false
	} else {
		return reflekt.AsBool(this.Values[0])
	}
}

// Bools returns values as bool array (values will be false, if not parsable to float64)
func (this *parameter) Bools() []bool {
	if this.Values == nil {
		return nil
	} else {
		res := make([]bool, this.Count())
		for i, v := range this.Values {
			res[i] = reflekt.AsBool(v)
		}
		return res
	}
}

// Time is a date time representation of the value with a provided format.
// If no format is provided, then `2006-01-02 15:04:05` will be used
func (this *parameter) Time(format ...string) (*time.Time, error) {
	if this.Values == nil {
		return nil, nil
	} else {
		f := "2006-01-02 15:04:05"
		if len(format) > 0 {
			f = format[0]
		}
		if t, err := time.Parse(f, this.Values[0]); err != nil {
			return nil, err
		} else {
			return &t, nil
		}
	}
}

// Times returns array for `time.Time` values, parsed from provided format.
// See `Time()`.
func (this *parameter) Times(format ...string) ([]time.Time, error) {
	if this.Values == nil {
		return nil, nil
	} else {
		f := "2006-01-02 15:04:05"
		if len(format) > 0 {
			f = format[0]
		}
		tt := make([]time.Time, len(this.Values))
		for i, v := range this.Values {
			if t, err := time.Parse(f, v); err != nil {
				return nil, err
			} else {
				tt[i] = t
			}
		}
		return tt, nil
	}
}

// Json assumes the input is a JSON string and parses into a standard map[string]interface{}
// Returns error, if not parsable (or eg array JSON).
//
// Helpful to allow complex inputs: `my-app do --foo '{"bar": "baz"}'`
func (this *parameter) Json() (map[string]interface{}, error) {
	if this.Values == nil {
		return nil, nil
	} else {
		m := make(map[string]interface{})
		if err := json.Unmarshal([]byte(this.Values[0]), &m); err != nil {
			return nil, err
		} else {
			return m, nil
		}
	}
}

// Jsons returns values as individual JSON strings. See `Json()` above.
func (this *parameter) Jsons() ([]map[string]interface{}, error) {
	if this.Values == nil {
		return nil, nil
	} else {
		res := make([]map[string]interface{}, len(this.Values))
		for i, v := range this.Values {
			m := make(map[string]interface{})
			if err := json.Unmarshal([]byte(v), &m); err != nil {
				return nil, err
			} else {
				res[i] = m
			}
		}
		return res, nil
	}
}
