package cli

import (
	"fmt"
	"os"
	"reflect"
)

// Clee is a command line interface object
type Cli struct {
	Name         string
	Version      string
	Description  string
	Commands     map[string]*Command
	Describer    Describer
	Registry     *Registry
	ErrorHandler func(error)
	Output       Output
}

func New(name, version, desc string) *Cli {
	return (&Cli{
		Name:        name,
		Version:     version,
		Description: desc,
		Commands:    make(map[string]*Command),
		Describer:   new(DefaultDescriber),
		Registry:    NewRegistry(),
		ErrorHandler: func(err error) {
			panic(err.Error())
		},
		Output: NewIoOutput(os.Stdout),
	}).Add(NewHelpCommand())
}

// Add is a builder method for adding a new command
func (this *Cli) Add(cmd *Command) *Cli {
	this.Commands[cmd.Name] = cmd.SetCli(this)
	return this
}

// New creates and adds a new command
func (this *Cli) New(name, usage string, call CallMethod) *Cli {
	return this.Add(NewCommand(name, usage, call).SetCli(this))
}

func (this *Cli) Register(v interface{}) *Cli {
	this.Registry.Register(v)
	return this
}

func (this *Cli) Run() {
	this.RunWith(os.Args[1:])
}

func (this *Cli) SetErrorHandler(handler func(error)) *Cli {
	this.ErrorHandler = handler
	return this
}

func (this *Cli) SetOutput(out Output) *Cli {
	this.Output = out
	return this
}

// Run the clee and be happy
func (this *Cli) RunWith(args []string) {
	if len(args) < 1 {
		fmt.Print(this.Describer.Cli(this))
	} else if c, ok := this.Commands[args[0]]; ok {
		method := c.Call.Type()
		input := make([]reflect.Value, method.NumIn())
		for i := 0; i < method.NumIn(); i++ {
			t := method.In(i)
			s := t.String()
			if s == "*cli.Cli" {
				input[i] = reflect.ValueOf(this)
			} else if s == "*cli.Command" {
				input[i] = reflect.ValueOf(c)
			} else if this.Registry.Has(s) {
				input[i] = this.Registry.Get(s)
			} else {
				this.ErrorHandler(fmt.Errorf("Missing parameter %s", s))
				return
			}
		}

		if err := c.Parse(args[1:]); err != nil {
			this.ErrorHandler(fmt.Errorf("Parse error: %s", err))
			return
		}

		res := c.Call.Call(input)

		if len(res) > 0 && res[0].Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) && !res[0].IsNil() {
			this.ErrorHandler(res[0].Interface().(error))
		}
	} else {
		this.ErrorHandler(fmt.Errorf("Command \"%s\" unknown", args[0]))
	}
}

// SetDescriber is builder method and switches describer
func (this *Cli) SetDescriber(desc Describer) *Cli {
	this.Describer = desc
	return this
}
