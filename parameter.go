package cli

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

type ParserMethod func(name, value string) (string, error)
type ValidatorMethod func(name, value string) error

// parameter is core for Argument and Option
type parameter struct {

	// Name is used for describing and accessing this argument
	Name string

	// Usage is a short description of this argument
	Usage string

	// Descriptions is a lengthy elaboration of the purpose, use-case, life-story of this argument
	Description string

	// Required determines whether command can execute without this argument
	// Should NOT be changed after adding with `AddCommand` from `Command`
	Required bool

	// Whether multiple values are allowed. Only
	Multiple bool

	// Default is used if no value is provided
	Default string

	// Value holds what was provided on the command line
	Values []string

	// Parser generates value to be returned on calling `Parsed()` (optional)
	Parser ParserMethod

	// Regex for checking if input value can be accepted
	Regex *regexp.Regexp

	// Validator method to check if input value can be accepted
	Validator ValidatorMethod
}

/*
Arguments must be provided immediately after the command in the order they were added.
Non required arguments must be orderd after required arguments. Only one argument is allowed
to contain multiple values and it needs to be the last one.
*/
type Argument struct {
	parameter
}

func NewArgument(name, usage, _default string, required, multiple bool) *Argument {
	return &Argument{
		parameter: parameter{
			Name:     name,
			Usage:    usage,
			Required: required,
			Multiple: multiple,
			Default:  _default,
		},
	}
}

/*func (this *Argument) Format(templ string) string {
	res := strings.Replace(templ, ":name:", this.Name, -1)
	res = strings.Replace(res, ":usage:", this.Usage, -1)
	res = strings.Replace(res, ":description:", this.Description, -1)
	res = strings.Replace(res, ":required:", fmt.Sprintf("%v", this.Required), -1)
	res = strings.Replace(res, ":multiple:", fmt.Sprintf("%v", this.Multiple), -1)
	res = strings.Replace(res, ":default:", this.Default, -1)
	return res
}*/

// Option is like an argument, but with a single or double dash in front
type Option struct {
	parameter

	// Alias can be a shorter name
	Alias string

	// If is a flag, then no value can be assigned (if present, then bool true)
	Flag bool
}

func NewOption(name, alias, usage, _default string, required, multiple bool) *Option {
	return &Option{
		parameter: parameter{
			Name:     name,
			Usage:    usage,
			Required: required,
			Multiple: multiple,
			Default:  _default,
		},
		Alias: alias,
	}
}

func (this *Option) IsFlag() *Option {
	this.Flag = true
	return this
}

/*
---------------------
Builder
---------------------
*/

// SetUsage is a builder method to set usage
func (this *parameter) SetUsage(v string) *parameter {
	this.Usage = v
	return this
}

// SetDescription is a builder method to set description
func (this *parameter) SetDescription(v string) *parameter {
	this.Description = v
	return this
}

// SetDefault is a builder method to set default value
func (this *parameter) SetDefault(v string) *parameter {
	this.Default = v
	return this
}

// SetDefault is a builder method to set default value
func (this *parameter) SetParser(v ParserMethod) *parameter {
	this.Parser = v
	return this
}

// SetDefault is a builder method to set default value
func (this *parameter) SetRegex(r *regexp.Regexp) *parameter {
	this.Regex = r
	return this
}

// SetDefault is a builder method to set default value
func (this *parameter) SetValidator(v ValidatorMethod) *parameter {
	this.Validator = v
	return this
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
		if this.Validator != nil {
			if err := this.Validator(this.Name, val); err != nil {
				return fmt.Errorf(print(err.Error()))
			}
		}
		if this.Parser != nil {
			if p, err := this.Parser(this.Name, val); err != nil {
				return fmt.Errorf(print(err.Error()))
			} else {
				val = p
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

// Strings representation of the value (can be nil)
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

// Ints representation of the value (can be nil, if non given or slice of 0-values for)
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

// Float representation of the value (will be 0, if not given or not parsable)
func (this *parameter) Float() float64 {
	if this.Values == nil {
		return 0
	} else {
		return reflekt.AsFloat(this.Values[0])
	}
}

// Ints representation of the value (can be nil, if non given or slice of 0-values for)
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

// Float representation of the value (will be false, if not given or not parsable)
func (this *parameter) Bool() bool {
	if this.Values == nil {
		return false
	} else {
		return reflekt.AsBool(this.Values[0])
	}
}

// Ints representation of the value (can be nil, if non given or slice of 0-values for)
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

// Date is a date time representation of the value
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
// Returns error, if not parsable (or eg array JSON)
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
