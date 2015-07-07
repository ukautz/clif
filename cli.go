package clif

import (
	"os"
	"reflect"
	"fmt"
	"strings"
	"os/signal"
)

// HeraldCallback is type for creator method, which can be registered with
// the `Herald()` method.
type HeraldCallback func(*Cli) *Command

// NamedParameters ...
type NamedParameters map[string]interface{}

// Cli is a command line interface object
type Cli struct {

	// Name is the name of the console application used in the generated help
	Name string

	// Version is used in the generated help
	Version string

	// Description is used in the generated help
	Description string

	// Commands contain all registered commands and can be manipulated directly
	Commands map[string]*Command

	// Heralds contain list of command-create-callbacks which will be executed on `Run()`
	Heralds []HeraldCallback

	// Registry is a container holding objects for injection
	Registry *Registry

	// DefaultOptions contains options which are added to all commands early
	// in the `Run()` call.
	DefaultOptions []*Option

	// OnInterrupt, when set with `SetOnInterrupt`, is callback which is executed
	// if user triggers interrupt (ctrl+c). If an error is returned, then the
	// cli application will die with a non-zero status and print the error message.
	onInterrupt func() error

	// interruptChan is channel for set interrupt callback
	interruptChan chan os.Signal
}

// New constructs new cli
func New(name, version, desc string) *Cli {
	this := &Cli{
		Name:           name,
		Version:        version,
		Description:    desc,
		Commands:       make(map[string]*Command),
		Heralds:        make([]HeraldCallback, 0),
		Registry:       NewRegistry(),
		DefaultOptions: make([]*Option, 0),
	}
	this.Add(NewHelpCommand())
	out := NewColorOutput(os.Stdout)
	this.Register(this).SetOutput(out).SetInput(NewDefaultInput(os.Stdin, out))
	return this
}

// Add is a builder method for adding a new command
func (this *Cli) Add(cmd ...*Command) *Cli {
	for _, c := range cmd {
		this.Commands[c.Name] = c.SetCli(this)
	}
	return this
}

// NewDefaultOption creates and adds a new option to default list.
func (this *Cli) NewDefaultOption(name, alias, usage, _default string, required, multiple bool) *Cli {
	return this.AddDefaultOptions(NewOption(name, alias, usage, _default, required, multiple))
}

// AddDefaultOptions adds a list of options to default options.
func (this *Cli) AddDefaultOptions(opts ...*Option) *Cli {
	this.DefaultOptions = append(this.DefaultOptions, opts...)
	return this
}

// Herald registers command constructors, which will be executed in `Run()`.
func (this *Cli) Herald(cmd ...HeraldCallback) *Cli {
	for _, c := range cmd {
		this.Heralds = append(this.Heralds, c)
	}
	return this
}

// New creates and adds a new command
func (this *Cli) New(name, usage string, call CallMethod) *Cli {
	return this.Add(NewCommand(name, usage, call).SetCli(this))
}

// Output is shorthand for currently registered output
func (this *Cli) Output() Output {
	t := reflect.TypeOf((*Output)(nil)).Elem()
	out := this.Registry.Get(t.String())
	return out.Interface().(Output)
}

// RegisterAs is builder method and registers object in registry
func (this *Cli) Register(v interface{}) *Cli {
	this.Registry.Register(v)
	return this
}

// RegisterAs is builder method and registers object under alias in registry
func (this *Cli) RegisterAs(n string, v interface{}) *Cli {
	this.Registry.Alias(n, v)
	return this
}

// RegisterNamed registers a parameter for injecting in a named map[string]inteface{}
func (this *Cli) RegisterNamed(n string, v interface{}) *Cli {
	this.Registry.Alias(fmt.Sprintf("N:%s", n), v)
	return this
}

// Run with OS command line arguments
func (this *Cli) Run() {
	this.RunWith(os.Args[1:])
}

// RunWith runs the cli with custom list of arguments
func (this *Cli) RunWith(args []string) {
	for _, cb := range this.Heralds {
		this.Add(cb(this))
	}
	this.Heralds = make([]HeraldCallback, 0)
	for _, cmd := range this.Commands {
		for _, opt := range this.DefaultOptions {
			cmd.AddOption(opt)
		}
	}
	name := ""
	cargs := []string{}
	if len(args) < 1 || (len(args) == 1 && (args[0] == "-h" || args[0] == "--help")) {
		this.Output().Printf(DescribeCli(this))
		return
	} else {
		name = args[0]
		cargs = args[1:]
	}

	if c, ok := this.Commands[name]; ok {

		// parse arguments & options
		err := c.Parse(cargs)
		if help := c.Option("help"); help != nil && help.Bool() {
			this.Output().Printf(DescribeCommand(c))
			return
		}
		if err != nil {
			this.Output().Printf(DescribeCommand(c))
			Die("Parse error: %s", err)
		}

		// execute callback & handle result
		if res, err := this.Call(c); err != nil {
			Die(err.Error())
		} else {
			errType := reflect.TypeOf((*error)(nil)).Elem()
			if len(res) > 0 && res[0].Type().Implements(errType) && !res[0].IsNil() {
				Die("Failure in execution: %s", res[0].Interface().(error))
			}
		}
	} else {
		Die("Command \"%s\" unknown", args[0])
	}
}

// Call executes command by building all input parameters based on objects
// registered in the container and running the callback.
func (this *Cli) Call(c *Command) ([]reflect.Value, error) {

	// build callback arguments and execute
	this.Register(c)
	method := c.Call.Type()
	input := make([]reflect.Value, method.NumIn())
	named := NamedParameters(make(map[string]interface{}))
	namedType := reflect.TypeOf(named).String()
	namedIndex := -1
	for i := 0; i < method.NumIn(); i++ {
		t := method.In(i)
		s := t.String()
		if this.Registry.Has(s) {
			input[i] = this.Registry.Get(s)
		} else if s == namedType {
			if namedIndex > -1 {
				return nil, fmt.Errorf("Callback has more than the one allowed input parameter of type %s, which is used to inject named parameters", namedType)
			}
			namedIndex = i
		} else {
			return nil, fmt.Errorf("Missing callback parameter %s", s)
		}
	}
	if namedIndex > -1 {
		this.Registry.Reduce(func (name string, value interface{}) bool {
			if strings.Index(name, "N:") == 0 {
				named[name[2:]] = value
			}
			return false
		})
		input[namedIndex] = reflect.ValueOf(NamedParameters(named))
	}

	return c.Call.Call(input), nil
}

// SetDescription is builder method and sets description
func (this *Cli) SetDescription(v string) *Cli {
	this.Description = v
	return this
}

// SetOutput is builder method and replaces current input
func (this *Cli) SetInput(in Input) *Cli {
	t := reflect.TypeOf((*Input)(nil)).Elem()
	this.Registry.Alias(t.String(), in)
	return this
}

// SetOutput is builder method and replaces current output
func (this *Cli) SetOutput(out Output) *Cli {
	t := reflect.TypeOf((*Output)(nil)).Elem()
	this.Registry.Alias(t.String(), out)
	return this
}

// SetOnInterrupt sets callback for interrupt and associates signal
func (this *Cli) SetOnInterrupt(cb func() error) *Cli {
	this.onInterrupt = cb

	if this.interruptChan == nil {
		this.interruptChan = make(chan os.Signal, 1)
		signal.Notify(this.interruptChan, os.Interrupt)
		go func() {
			<-this.interruptChan
			if err := this.onInterrupt(); err != nil {
				Die(err.Error())
			} else {
				Exit(0)
			}
		}()
	}

	return this
}