package clif

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var (
	rxFlag = regexp.MustCompile(`^(?:true|false|yes|no)$`)
)

// CallMethod is the interface for functions used as command callbacks.
// `interface{}` is used so that those functions can have arbitrary input and output signatures.
// However, they still must be functions.
type CallMethod interface{}

// Command represents a named callback with a set of arguments and options
type Command struct {

	// Cli back-references the Cli in which the command is registered
	Cli *Cli

	// Name is the unique (within Cli scope) call-name of the command
	Name string

	// Usage is a shorthand description of what the command does. Used in help output.
	Usage string

	// Description is a long elaboration on what the command does. Used in help output.
	Description string

	// Options contain all the registered options of the command.
	Options []*Option

	// Arguments contain all the registered arguments of the command.
	Arguments []*Argument

	// Call holds reflections of the callback.
	Call reflect.Value

	// PreCall is optional method which will be executed before command Call
	preCall *reflect.Value

	// PostCall is optional method which will be executed after command Call
	postCall *reflect.Value
}

// DefaultHelpOption is the "--help" option, which is (per default) added to any command.
var DefaultHelpOption = &Option{
	parameter: parameter{
		Name:        "help",
		Usage:       "Display this help message",
		Description: "Display this help message",
	},
	Alias: "h",
	Flag:  true,
}

// DefaultOptions are prepended to any newly created command. Will be added
// immediately in the `NewCommand` call. See also `cli.NewDefaultOption()` and
// `cli.AddDefaultOptions()`.
var DefaultOptions = []*Option{DefaultHelpOption}

// NewCommand constructs a new command
func NewCommand(name, usage string, call CallMethod) *Command {
	ref := reflect.ValueOf(call)
	if ref.Kind() != reflect.Func {
		panic(fmt.Sprintf("Call must be method, but is %s", ref.Kind()))
	}
	return &Command{
		Name:      name,
		Usage:     usage,
		Options:   DefaultOptions,
		Arguments: make([]*Argument, 0),
		Call:      ref,
	}
}

// SetCli is builder method and sets the Cli back-reference
func (this *Command) SetCli(c *Cli) *Command {
	this.Cli = c
	return this
}

// SetDescription is builder method setting description
func (this *Command) SetDescription(desc string) *Command {
	this.Description = desc
	return this
}

func (this *Command) SetPreCall(call CallMethod) *Command {
	ref := reflect.ValueOf(call)
	if ref.Kind() != reflect.Func {
		panic(fmt.Sprintf("PreCall must be method, but is %s", ref.Kind()))
	}
	this.preCall = &ref

	return this
}

func (this *Command) SetPostCall(call CallMethod) *Command {
	ref := reflect.ValueOf(call)
	if ref.Kind() != reflect.Func {
		panic(fmt.Sprintf("PostCall must be method, but is %s", ref.Kind()))
	}
	this.postCall = &ref

	return this
}

// Parse extracts options and arguments from command line arguments
func (this *Command) Parse(args []string) error {
	argNum := 0
	var lastArg *Argument
	argc := len(args)
	la := len(this.Arguments)
	for i := 0; i < argc; i++ {
		a := args[i]
		if len(a) > 0 && a[0:1] == "-" {
			n := strings.TrimLeft(a, "-")
			v := ""
			vv := false

			// is --opt=foo format
			if p := strings.Index(n, "="); p > -1 {
				vv = true
				if p == 0 {
					return fmt.Errorf("Malformed option \"%s\"", a)
				}
				pp := strings.SplitN(n, "=", 2)
				n = pp[0]
				if len(pp) == 2 {
					v = pp[1]
				}
			}
			o := this.Option(n)
			if o == nil {
				return fmt.Errorf("Unrecognized option \"%s\"", a)
			}

			// not a flag: must have value
			// is a flag: must not have non-flag value
			if !o.Flag && !vv {
				if i+1 == argc || args[i+1][0:1] == "-" {
					return fmt.Errorf("Missing value for option \"%s\"", a)
				} else {
					i++
					v = args[i]
				}
			} else if o.Flag && vv && !rxFlag.MatchString(v) {
				return fmt.Errorf("Flag \"%s\" cannot have value", a)
			} else if o.Flag {
				v = "true"
			}
			if err := o.Assign(v); err != nil {
				return err
			}
		} else {
			if lastArg == nil || !lastArg.Multiple {
				if argNum+1 > la {
					return fmt.Errorf("Too many arguments. Expected (at most) %d, got %d", la, argNum+1)
				} else {
					lastArg = this.Arguments[argNum]
				}
			}
			if err := lastArg.Assign(a); err != nil {
				return err
			}
			argNum++
		}
	}

	for _, a := range this.Arguments {
		if err := this.postSetupParam(a, "Argument"); err != nil {
			return err
		}
	}
	for _, o := range this.Options {
		if err := this.postSetupParam(o, "Option"); err != nil {
			return err
		}
	}

	return nil
}

func (this *Command) postSetupParam(x interface{}, t string) error {
	var p *parameter
	if a, ok := x.(*Argument); ok {
		p = &(a.parameter)
	} else {
		p = &(x.(*Option).parameter)
	}

	if len(p.Values) == 0 {
		v := ""
		if p.Env != "" {
			v = os.Getenv(p.Env)
		}
		if v == "" && p.Default != "" {
			v = p.Default
		}
		if v != "" {
			if err := p.Assign(v); err != nil {
				return err
			}
		}
	}
	if p.Required && p.Count() == 0 {
		return fmt.Errorf("%s \"%s\" is required but missing", t, p.Name)
	} else {
		return nil
	}
}

// NewArgument is builder method to construct and add a new argument
func (this *Command) NewArgument(name, usage, _default string, required, multiple bool) *Command {
	return this.AddArgument(NewArgument(name, usage, _default, required, multiple))
}

// AddArgument is builder method to add a new argument
func (this *Command) AddArgument(v *Argument) *Command {
	var prev *Argument
	if l := len(this.Arguments); l > 0 {
		prev = this.Arguments[l-1]
	}
	if prev != nil {
		if v.Required && !prev.Required {
			panic("Cannot add required argument after optional argument")
		} else if prev.Multiple {
			panic("Cannot add argument after multiple style argument")
		}
	}
	if this.Argument(v.Name) != nil {
		panic(fmt.Sprintf("Argument with name \"%s\" already existing", v.Name))
	} else if this.Option(v.Name) != nil {
		panic(fmt.Sprintf("Option with name or alias \"%s\" already existing", v.Name))
	}
	this.Arguments = append(this.Arguments, v)

	return this
}

// NewFlag adds a new flag option
func (this *Command) NewFlag(name, alias, usage string, multiple bool) *Command {
	return this.AddOption(NewFlag(name, alias, usage, multiple))
}

// NewOption is builder method to construct and add a new option
func (this *Command) NewOption(name, alias, usage, _default string, required, multiple bool) *Command {
	return this.AddOption(NewOption(name, alias, usage, _default, required, multiple))
}

// AddOption is builder method to add a new option
func (this *Command) AddOption(v *Option) *Command {
	if this.Option(v.Name) != nil {
		panic(fmt.Sprintf("Option with name or alias \"%s\" already existing", v.Name))
	} else if this.Argument(v.Name) != nil {
		panic(fmt.Sprintf("Argument with name \"%s\" already existing", v.Name))
	} else if v.Alias != "" {
		if this.Option(v.Alias) != nil {
			panic(fmt.Sprintf("Option with name or alias \"%s\" already existing", v.Alias))
		} else if this.Argument(v.Alias) != nil {
			panic(fmt.Sprintf("Cannot use alias: Argument with name \"%s\" already existing", v.Alias))
		}
	}
	this.Options = append(this.Options, v)
	return this
}

// Argument provides access to registered, named arguments.
func (this *Command) Argument(name string) *Argument {
	for _, a := range this.Arguments {
		if a.Name == name {
			return a
		}
	}
	return nil
}

// Option provides access to registered, named options.
func (this *Command) Option(name string) *Option {
	for _, o := range this.Options {
		if o.Name == name || o.Alias == name {
			return o
		}
	}
	return nil
}

// Input returns map containing whole input values (of all options, all arguments)
func (this *Command) Input() map[string][]string {
	res := make(map[string][]string)
	for _, o := range this.Options {
		if len(o.Values) > 0 {
			res[o.Name] = o.Values
		}
	}
	for _, a := range this.Arguments {
		if len(a.Values) > 0 {
			res[a.Name] = a.Values
		}
	}
	return res
}
