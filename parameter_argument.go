package clif
import "regexp"

// Arguments must be provided immediately after the command in the order they were
// added. Non required arguments must be ordered after required arguments. Only
// one argument is allowed to contain multiple values and it needs to be the last one.
type Argument struct {
	parameter
}

// NewArgument constructs a new argument
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

// SetUsage is builder method to set the usage description. Usage is a short
// account of what the argument is used for, for help generaiton.
func (this *Argument) SetUsage(v string) *Argument {
	this.Usage = v
	return this
}

// SetDescription is a builder method to sets argument description. Description
// is an elaborate explanation which is used in help generation.
func (this *Argument) SetDescription(v string) *Argument {
	this.Description = v
	return this
}

// SetDefault is a builder method to set default value. Default value is used
// if the argument is not provided (after environment variable).
func (this *Argument) SetDefault(v string) *Argument {
	this.Default = v
	return this
}

// SetEnv is a builder method to set environment variable name, from which to
// take the value, if not provided (before default).
func (this *Argument) SetEnv(v string) *Argument {
	this.Env = v
	return this
}

// SetParse is a builder method to set setup call on value. The setup call
// must return a replace value (can be unchanged) or an error, which stops
// evaluation and returns error to user.
func (this *Argument) SetParse(v ParseMethod) *Argument {
	this.Parse = v
	return this
}

// SetRegex is a builder method to set regular expression which is used to
// check the the argument input (in case of multiple: each will be checked)
func (this *Argument) SetRegex(r *regexp.Regexp) *Argument {
	this.Regex = r
	return this
}
