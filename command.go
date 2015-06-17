package clif

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	rxFlag = regexp.MustCompile(`^(?:true|false|yes|no)$`)
)

type CallMethod interface{}

// Command ..
type Command struct {
	Cli         *Cli
	Name        string
	Usage       string
	Description string
	Options     []*Option
	Arguments   []*Argument
	Call        reflect.Value
}

// DefaultOptions are prepended to any newly created command
var DefaultHelpOption = &Option{
	parameter: parameter{
		Name:        "help",
		Usage:       "Display this help message",
		Description: "Display this help message",
	},
	Alias: "h",
	Flag:  true,
}

// DefaultOptions are prepended to any newly created command
var DefaultOptions = []*Option{DefaultHelpOption}

// NewCommand constructs new command
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

// SetDescription is builder method setting desciption
func (this *Command) SetCli(c *Cli) *Command {
	this.Cli = c
	return this
}

// SetDescription is builder method setting desciption
func (this *Command) SetDescription(desc string) *Command {
	this.Description = desc
	return this
}

// Parse command line args to argument
func (this *Command) Parse(args []string) error {
	argNum := 0
	var lastArg *Argument
	argc := len(args)
	la := len(this.Arguments)
	for i := 0; i < argc; i++ {
		a := args[i]
		//l := len(a)
		if a[0:1] == "-" {
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
		if len(a.Values) == 0 && a.Default != "" {
			a.Values = []string{a.Default}
		}
		if a.Required && len(a.Values) == 0 {
			return fmt.Errorf("Argument \"%s\" is required but missing", a.Name)
		}
	}
	for _, o := range this.Options {
		if len(o.Values) == 0 && o.Default != "" {
			o.Values = []string{o.Default}
		}
		if o.Required && len(o.Values) == 0 {
			return fmt.Errorf("Option \"%s\" is required but missing", o.Name)
		}
	}

	return nil
}

// AddArgument is builder method to add a new argument
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

// AddOption is builder method to add a new option
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

func (this *Command) Argument(name string) *Argument {
	for _, a := range this.Arguments {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func (this *Command) Option(name string) *Option {
	for _, o := range this.Options {
		if o.Name == name || o.Alias == name {
			return o
		}
	}
	return nil
}

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
