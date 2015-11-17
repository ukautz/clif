![status](https://api.travis-ci.org/ukautz/clif.svg?branch=v0)
[![GoDoc](https://godoc.org/gopkg.in/ukautz/clif.v0?status.svg)](http://godoc.org/gopkg.in/ukautz/clif.v0)


CLIF - Command line interface framework
=======================================

Go framework for rapid command line application development.

Example
-------

![demo](https://cloud.githubusercontent.com/assets/600604/8886731/c0a517e0-326f-11e5-8349-6ebee2cb8de5.gif)

```go
package main

import "gopkg.in/ukautz/clif.v0"

func main() {
	clif.New("My App", "1.0.0", "An example application").
		New("hello", "The obligatory hello world", func(out clif.Output) {
			out.Printf("Hello World\n")
		}).
		Run()
}
```

- - -

* [Install](#install)
* [Getting started](#getting-started)
* [Commands](#commands)
  * [Callback functions](#callback-functions)
    * [Named](#named)
    * [Default objects](#default-objects)
  * [Arguments and Options](#arguments-and-options)
    * [Arguments](#arguments)
    * [Options](#options)
      * [Flags](#flags)
    * [Validation &amp; (Parsing | Transformation)](#validation--parsing--transformation)
    * [Environment variables &amp; default](#environment-variables--default)
    * [Default options](#default-options)
* [Input &amp; Output](#input--output)
  * [Input](#input)
    * [Ask &amp; AskRegex](#ask--askregex)
    * [Confirm](#confirm)
    * [Choose](#choose)
  * [Output &amp; formatting](#output--formatting)
    * [Outputs](#outputs)
    * [Styles](#styles)
* [Real-life example](#real-life-example)
* [See also](#see-also)

- - -

Install
-------

``` bash
$ go get gopkg.in/ukautz/clif.v0
```

Getting started
---------------

One the one side, CLIF's *builder*-like API can be easily used for rapid development of small, single purpose tools. On the other side, CLIF is designed with complex console applications in mind.

Commands
--------

Commands must have a unique name and can have additional arguments and options.

``` go
cmd1 := clif.NewCommand("name", "A description", callBackFunction)
cmd2 := clif.NewCommand("other", "Another description", callBackFunction2)
```

The `name` is used from the command line to call the command:

```bash
$ ./app name
$ ./app other
```

### Callback functions

Callback functions can have arbitrary parameters. CLIF uses a small, built-in (signatur) injection container which allows you to register any kind of object (`struct` or `interface`) beforehand.

So you can register any object (interface{}, struct{} .. and anything else, see [below](#named)) in your bootstrap and then "require" those instances by simply putting them in the command callback signature:

```go
// Some type definition
type MyFoo struct {
    X int
}

func main() {
    // init cli
    cli := clif.New("My App", "1.0.0", "An example application")

    // register object instance with container
    foo := &MyFoo{X: 123}
    cli.Register(foo)

    // Create command with callback using the peviously registered instance
    cli.NewCommand("foo", "Call foo", func (foo *MyFoo) {
        // do something with foo
    })

    cli.Run()
}
```

Using interfaces is possible as well, but a bit less elegant:

```go
// Some interface
type MyBar interface {
    Bar() string
}

// Some type
type MyFoo struct {
}

// implement interface
func (m *MyFoo) Bar() string {
    return "bar"
}

func main() {
    // init cli
    cli := clif.New("My App", "1.0.0", "An example application")

    // create object, which implements MyBar:
    foo := &MyFoo{}
    t := reflect.TypeOf((*MyBar)(nil)).Elem()
    cli.RegisterAs(t.String(), foo)

    // Register command with callback using the type
    cli.NewCommand("bar", "Call bar", func (bar MyBar) {
        // do something with bar
    })

    cli.Run()
}
```

#### Named

Everything works great if you only have a single instance of any object of a specific type.
However, if you need more than one instance (which might often be the case for primitive
types, such as `int` or `string`) you can use named registering:

```go
// Register abitrary objects under unique name
cli.RegisterNamed("foo", new(MyFoo)).
    RegisterNamed("bar", 123).
    RegisterNamed("baz", "bla")

// Register command with callback named container
cli.NewCommand("bar", "Call bar", func (named clif.NamedParameters) {
    asMap := map[string]interface{}(named)
    fmt.Println(asMap["baz"].(string))
})
```

**Note**: If you want to use the named feature, you cannot `Register()` any `NamedParameters`
instance yourself, since "normally" registered objects are evaluated before named.

#### Default objects

CLIF pre-populates the dependency container with a couple of built-in objects:

* The `Output` (formatted output helper, see below), eg `func (out clif.Output) { .. }`
* The `Input` (input helper, see below), eg `func (in clif.Input) { .. }`
* The `*Cli` instance itself, eg `func (c *clif.Cli) { .. }`
* The current `*Command` instance, eg `func (o *clif.Command) { .. }`

### Arguments and Options

CLIF can deal with arguments and options. The difference being:

* **Arguments** come after the command name. They are identified by their position.
* **Options** have no fixed position. They are identified by their `--opt-name` (or alias, eg `-O`)

Of course you can use arguments and options at the same time..

#### Arguments

Arguments are additional command line parameters which come after the command name itself.

``` go
cmd := clif.NewCommand("hello", "A description", callBackFunction)
	.NewArgument("name", "Name for greeting", "", true, false)

arg := cmd.NewAgument("other", "Something ..", "default", false, true)
cmd.AddArgument(arg)
```

Arguments consist of a *name*, a *description*, an optional *default* value a *required* flag and a *multiple* flag.

``` bash
$ ./my-app hello the-name other1 other2 other3
#            ^      ^       ^       ^     ^
#            |      |       |       |     |
#            |      |       |       | third "other" arg
#            |      |       |  second "other" arg
#            |      |  first "other" arg
#            |  the "name" arg
#        command name
```

Position of arguments matters. Make sure you add them in the right order. And: **required** arguments must come before optional arguments (makes sense, right?). There can be only one **multiple** argument at all and, of course, it must be the last (think: variadic).

You can access the arguments by injecting the command instance `*clif.Command` into the callback and calling the `Argument()` method. You can choose to interpret the argument as `String()`, `Int()`, `Float()`, `Bool()`, `Time()` or `Json()`. Multiple arguments can be accessed with `Strings()`, `Ints()` .. and so on. `Count()` gives the amount of (provided) multiple arguments and `Provided()` returns bool for optional arguments. Please see [parameter.go](parameter.go) for more.

``` go
func callbackFunctionI(c *clif.Command) {
	// a single
	name := c.Argument("name").String()

	// a multiple
	others := c.Argument("other").Strings()

	// .. do something ..
}
```

#### Options

Options have no fixed position, meaning `./app --foo --bar` and `./app --bar --foo` are equivalent. Options are referenced by their name (eg `--name`) or alias (eg `-n`). Unless the option is a flag (see below) it must have a value. The value must immediately follow the option. Valid forms are: `--name value`, `--name=value`, `-n value` and `-n=value`.

Options must come before the command, unless they use the `=` separator. For example: `./app command --opt value` is valid, `./app --opt=value command` is valid but `./app --opt value command` is not valid (since it becomes impossible to distinguish between command and value).

``` go
cmd := clif.NewCommand("hello", "A description", callBackFunction)
	.NewOption("name", "n", "Name for greeting", "", true, false)

arg := cmd.NewOption("other", "O", "Something ..", "default", false, true)
cmd.AddOption(arg)
```

Now:

``` bash
$ ./my-app hello --other bar -n Me -O foo
#                       ^       ^    ^
#                       |       |    |
#                       |       |  second other opt with value
#                       |   name opt with value
#                  first other opt with value
```

You can access options the same way as arguments, just use `Option()` instead.

``` go
func callbackFunctionI(c *clif.Command) {
	name := c.Option("name").String()
	others := c.Option("other").Strings()
	// .. do something ..
}
```

##### Flags

There is a special kind of option, which does not expect a parameter: the flag. As options, their position is arbitrary.

``` go
flag := clif.NewOption("my-flag", "f", "Something ..", "", false, false).IsFlag()
cmd := clif.NewCommand("hello", "A description", callBackFunction).AddOption(flag)
```

When using the option, you dont need to (nor can you) provide an argument:

```bash
$ ./my-app hello --my-flag
```

You want to use `Bool()` to check if a flag is provided:

``` go
func callbackFunctionI(c *clif.Command) {
	if c.Option("my-flag").Bool() {
		// ..
	}
}
```

#### Validation & (Parsing | Transformation)

You can validate/parse/transform the input using the `Parse` attribute of options or arguments. It can be (later on)
set using the `SetParse()` method:

``` go
// Validation example
arg := clif.NewArgument("my-int", "An integer", "", true, false).
    SetParse(func(name, value string) (string, error) {
        if _, err := strconv.Atoi(value); err != nil {
            return "", fmt.Errorf("Oops: %s is not an integer: %s", name, err)
        } else {
            return value, nil
        }
    })

// Transformation example
opt := clif.NewOption("client-id", "c", "The client ID", "", true, false).
    SetParse(func(name, value string) (string, error) {
        if strings.Index(value, "#") != 0 {
            return fmt.Sprintf("#%s", value), nil
        } else {
            return value, nil
        }
    })
```

There are a few built-in validators you can use out of the box:

* `clif.IsInt` - Checks for integer, eg `clif.NewOption(..).SetParse(clif.IsInt)`
* `clif.IsFloat` - Checks for float, eg `clif.NewOption(..).SetParse(clif.IsFloat)`

See [validators.go](validators.go).

#### Environment variables & default

The argument and option constructors (`NewArgument`, `NewOption`) already allow you to set a default. In addition you can set
the name of an environment variable, which will be used, if the parameter is not provided.

``` go
opt := clif.NewOption("client-id", "c", "The client ID", "", true, false).SetEnv("CLIENT_ID")
```

The order is:

1. Provided, eg `--config /path/to/config`
2. Environment variable, eg `CONFIG_FILE`
3. Default value, as provided in constructor or set via `SetDefault()`

**Note**: A *required* parameter must have a value, but it does not care whether it came from input, via environment variable or as a default value.

#### Default options

Often you need one or multiple options on every or most commands. The usual `--verbose` or `--config /path..` are common examples.
CLIF provides two ways to deal with those.

1. Modifying/extending `clif.DefaultOptions` (it's pre-filled with the `--help` option, which is `clif.DefaultHelpOption`)
2. Calling `AddDefaultOptions()` or `NewDefaultOption()` on an instance of `clif.Cli`

The former is global (for any instance of `clif.Cli`) and assigned to any new command (created by the `NewCommand` constructor). The latter is applied when `Run()` is called and is in the scope of a single `clif.Cli` instance.

**Note**: A helpful patterns is combining default options and the injection container/registry. Following an example parsing a config file, which can be set on any command with `--config /path..` or as an environment variable and has a default path.

```go

type Conf struct {
    Foo string
    Bar string
}

func() main {

    // init new cli app
    cli := clif.New("my-app", "1.2.3", "My app that does something")

    // register default option, which fills injection container with config instance
    configOpt := clif.NewOption("config", "c", "Path to config file", "/default/config/path.json", true, false).
        SetEnv("MY_APP_CONFIG").
        SetParse(function(name, value string) (string, error) {
            conf := new(Conf)
            if raw, err := ioutil.ReadFile(value); err != nil {
                return "", fmt.Errorf("Could not read config file %s: %s", value, err)
            } else if err = json.Unmarshal(raw, conf); err != nil {
                return "", fmt.Errorf("Could not unmarshal config file %s: %s", value, err)
            } else if conf.Foo == "" {
                return "", fmt.Errorf("Config %s is missing \"foo\"", value)
            } else {
                // register *Conf
                cli.Register(conf)
                return value, nil
            }
        })
    cli.AddDefaultOptions(configOpt)

    // Since *Conf was registered it can be used in any callback
    cli.New("anything", "Does anything", func(conf *Conf) {
        // do something with conf
    })

    cli.Run()
}
```

Input & Output
--------------

Of course, you can just use `fmt` and `os.Stdin`, but for convenience (and fancy output) there are `clif.Output` and `clif.Input`.

### Input

You can inject an instance of the `clif.Input` interface into your command callback. It provides small set of often used tools.

![input](https://cloud.githubusercontent.com/assets/600604/8886968/378a2668-3273-11e5-8bda-51b2b5cd127b.png)

#### Ask & AskRegex

Just ask the user a question then read & check the input. The question will be asked until the check/requirement is satisfied (or the user exits out with `ctrl+c`):

``` go
func callbackFunctionI(in clif.Input) {
	// Any input is OK
	foo := in.Ask("What is a foo", nil)

	// Validate input
	name := in.Ask("Who are you? ", func(v string) error {
		if len(v) > 0 {
			return nil
		} else {
			return fmt.Errorf("Didn't catch that")
		}
	})

	// Shorthand for regex validation
	count := in.AskRegex("How many? ", regexp.MustCompile(`^[0-9]+$`))

	// ..
}
```

*See `clif.RenderAskQuestion` for customization.*

#### Confirm

`Confirm()` ask the user a question until it is answered with `yes` (or `y`) or `no` (or `n`) and returns the response as `bool`.

``` go
func callbackFunctionI(in clif.Input) {
	if in.Confirm("Let's do it?") {
		// ..
	}
}
```

*See `clif.ConfirmRejection`, `clif.ConfirmYesRegex` and `clif.ConfirmNoRegex` for customization.*

#### Choose

`Choose()` is like a select in HTML and provides a list of options with descriptions to the user. The user then must choose (type in) one of the options. The choices will be presented to the user until a valid choice (one of the options) is provided.

``` go
func callbackFunctionI(in clif.Input) {
	father := in.Choose("Who is your father?", map[string]string{
		"yoda":  "The small, green guy",
		"darth": "The one with the smoker voice and the dark cape!",
		"obi":   "The old man with the light thingy",
	})

	if father == "darth" {
		// ..
	}
}
```

*See `clif.RenderChooseQuestion`, `clif.RenderChooseOption` and `clif.RenderChooseQuery` for customization.*

### Output & formatting

The `clif.Output` interface can be injected into any callback. It relies on a `clif.Formatter`, which does the actual formatting (eg colorizing) of the text.

#### Outputs

Per default, the `clif.DefaultInput` via `clif.NewColorOutput()` is used. It uses `clif.DefaultStyles`, which look like the screenshots you are seeing in this readme.

You can change the output like so:

``` go
cli := clif.New(..)
cli.SetOutput(clif.NewColorOutput().
    SetFormatter(clif.NewDefaultFormatter(clif.SunburnStyles))
```

#### Styles

Styles are applied by parsing (replacing) tokens like `<error>`, which would be substitude with `\033[31;1m` (using the default styles) resulting in a red coloring. Another example is `<reset>`, which is replaced with `\033[0m` leading to reset all colorings & styles.

There three built-in color styles (of course, you can extend them or add your own):

1. `DefaultStyles` - as you can see on this page
1. `SunburnStyles` - more yellow'ish
1. `WinterStyles` - more blue'ish


#### Helpers

Currently there are two output helpers available:

1. **Table renderer**: Render datasets in tables
2. **Progress bar**: Visualize arbitrary progress

##### Table


##### Progress bar

```go
func
pb := out.


Real-life example
-----------------

To provide you a usful'ish example, I've written a small CLI application called [repos](https://github.com/ukautz/repos).

See also
--------

There are a lot of [other approaches](https://github.com/avelino/awesome-go#command-line) you should have a look at.
