package clif

import "regexp"

// Option is a user input which is initialized with a single or double dash
// (eg "--foo" or "-f"). It may not be followed by a value, in which case it
// is considered a flag (see `IsFlag`). Options can be multiple inputs (no
// restrictions as there are for Arguments). An option can be required or optional.
// Options do not need to have any particular order.
type Option struct {
	parameter

	// Alias can be a shorter name
	Alias string

	// If is a flag, then no value can be assigned (if present, then bool true)
	Flag bool
}

// NewOption constructs new option
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

// NewFlag constructs new flag option
func NewFlag(name, alias, usage string, multiple bool) *Option {
	return &Option{
		parameter: parameter{
			Name:     name,
			Usage:    usage,
			Multiple: multiple,
		},
		Alias: alias,
		Flag:  true,
	}
}

// IsFlag marks an option as a flag. A Flag does not have any values. If it
// exists (eg "--verbose"), then it is automatically initialized with the string
// "true", which then can be checked with the `Bool()` method for actual `bool`
func (this *Option) IsFlag() *Option {
	this.Flag = true
	return this
}

// SetUsage is builder method to set the usage description. Usage is a short
// account of what the option is used for, for help generaiton.
func (this *Option) SetUsage(v string) *Option {
	this.Usage = v
	return this
}

// SetDescription is a builder method to sets option description. Description
// is an elaborate explanation which is used in help generation.
func (this *Option) SetDescription(v string) *Option {
	this.Description = v
	return this
}

// SetDefault is a builder method to set default value. Default value is used
// if the option is not provided (after environment variable).
func (this *Option) SetDefault(v string) *Option {
	this.Default = v
	return this
}

// SetEnv is a builder method to set environment variable name, from which to
// take the value, if not provided (before default).
func (this *Option) SetEnv(v string) *Option {
	this.Env = v
	return this
}

// SetParse is a builder method to set setup call on value. The setup call
// must return a replace value (can be unchanged) or an error, which stops
// evaluation and returns error to user.
func (this *Option) SetParse(v ParseMethod) *Option {
	this.Parse = v
	return this
}

// SetRegex is a builder method to set regular expression which is used to
// check the the option input (in case of multiple: each will be checked)
func (this *Option) SetRegex(r *regexp.Regexp) *Option {
	this.Regex = r
	return this
}
