![status](https://api.travis-ci.org/ukautz/clif.svg?branch=v0)

CLIF - Command line interface framework
=======================================

Go framework for rapid command line application development.

Example
-------

![extended](https://cloud.githubusercontent.com/assets/600604/8201524/50f04bdc-14d2-11e5-8ccb-591ef9e9b1f1.png)

```go
package main

import "gopkg.in/ukautz/clif.v0"

func main() {
	c := cli.New("My App", "1.0.0", "An example application").
		New("hello", "The obligatory hello world", func(out cli.Output) {
			out.Printf("Hello World\n")
		})
	c.Run()
}
```

- - -

* [Install](#install)
* [Getting started](#getting-started)
* [Commands](#commands)
  * [Callback functions](#callback-functions)
  * [Arguments and Options](#arguments-and-options)
    * [Arguments](#arguments)
    * [Options](#options)
      * [Flags](#flags)
    * [Validation](#validation)
* [Input &amp; Output](#input--output)
  * [Customizing Input](#customizing-input)
  * [Extending Output](#extending-output)
* [Patterns](#patterns)
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

```go
clif.New("My App", "1.0.0", "An example application").
	New("ls", "", func() { fmt.Printf("Foo\n") }).
	New("bar", "", func() { fmt.Printf("Bar\n") }).
	Run()
```

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

Callback functions can have arbitrary parameters. CLIF uses a small, built-in callback injection container which allows you to register any kind of object (`struct` or `interface`) beforehand.

```go
// Some type definition
type MyFoo struct {
    X int
}

// register object instance with container
foo := &MyFoo{X: 123}
cli.Register(foo)

// Register command with callback using the type
cli.NewCommand("foo", "Call foo", func (foo *MyFoo) {
    // do something with foo
})
```

Using interfaces is possible as well, but a bit less elegant:

```go
// Some type definition
type MyBar interface {
    Bar() string
}

// create object, which implements MyBar:
foo := &MyFoo{}
t := reflect.TypeOf((*MyBar)(nil)).Elem()
cli.RegisterAs(t.String(), foo)

// Register command with callback using the type
cli.NewCommand("bar", "Call bar", func (bar MyBar) {
    // do something with bar
})
```

#### Named

This works great if you only have a single instance of any object of a specific type.
However, if you need more than one instance (which might often be the case for primitive
types, such as `int` or `string`) you can use named registering:

```go
// Register abitrary objects under unique name
cli.RegisterNamed("foo", new(MyFoo)).
    RegisterNamed("bar", 123).
    RegisterNamed("baz", "bla")

// Register command with callback named container
cli.NewCommand("bar", "Call bar", func (named map[string]interface{}) {
    fmt.Println(named["baz"].(string))
})
```

**Note**: If you want to use the named feature, you cannot `Register()` any `map[string]interface{}`, since
"normally" registered objects are evaluated before named.

#### Default objects

CLIF pre-populates the dependency container with a couple of built-in objects:

* The `Output` (formatted output helper, see below)
* The `Input` (input helper, see below)
* The `*Cli` instance itself
* The current `*Command` instance

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

Position of arguments matters. Make sure you add them in the right order. And: **required** arguments must come before optional arguments (makes sense, right?). There can be only one **multiple** argument at all.

You can access the arguments by injecting the command instance into the callback and calling the `Argument()` method. You can choose to interpret the argument as `String()`, `Int()`, `Float()`, `Bool()`, `Time()` or `Json()`. Multiple arguments can be accessed with `Strings()`, `Ints()` .. and so on. `Count()` gives the amount of (provided) multiple arguments and `Provided()` returns bool for optional arguments. Pleas see [parameter.go](parameter.go) for more.

``` go
func callbackFunctionI(c *clif.Command) {
	name := c.Argument("name").String()
	others := c.Argument("other").Strings()
	// .. do something ..
}
```

#### Options

Options have no fixed position. They are referenced by their name (eg `--name`) or alias (eg `-n`).

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
#                       |       |  second other opt
#                       |   name opt
#                  first other opt
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

There is a special kind of option, which does not expect a parameter: the flag.

``` go
flag := clif.NewOption("my-flag", "f", "Something ..", "", false, false).IsFlag()
cmd := clif.NewCommand("hello", "A description", callBackFunction).AddOption(flag)
```

Usually, you want to use `Bool()` on flags:

``` go
func callbackFunctionI(c *clif.Command) {
	if c.Option("my-flag").Bool() {
		// ..
	}
}
```

#### Validation

You can provide argument or option validators, which are executed on parsing the command line input, before it is delegated to your callback.

``` go
arg := clif.NewArgument("my-int", "An integer", "", true, false)
arg.SetValidator(func(name, value string) error {
	if _, err := strconv.Atoi(value); err != nil {
		return fmt.Errorf("Oops: %s is not an integer: %s", name, err)
	} else {
		return nil
	}
})
```

There are a couple of built-in validators you can use out of the box:

* `clif.IsInt` - Checks for integer
* `clif.IsFloat` - Checks for float

See [validators.go](validators.go).

Input & Output
--------------

Of course, you can just use `fmt` and `os.Stdin`, but for convenience (and fancy output) there are `clif.Output` and `clif.Input`.

``` go
func callbackFunctionI(in clif.Input, out clif.Output) {
	name := in.Ask("Who are you? ", func(v string) error {
		if len(v) > 0 {
			return nil
		} else {
			return fmt.Errorf("Didn't catch that")
		}
	})
	father := in.Choose("Who is your father?", map[string]string{
		"yoda":  "The small, green guy",
		"darth": "NOOOOOOOO!",
		"obi":   "The old man with the light thingy",
	})

	out.Printf("Well, %s, ", name)
	if father != "darth" {
		out.Printf("<success>may the force be with you!<reset>\n")
	} else {
		out.Printf("<error>u bad!<reset>\n")
	}
}
```

![input](https://cloud.githubusercontent.com/assets/600604/8201525/510f730e-14d2-11e5-83aa-b238c804e98f.png)

### Customizing Input

There is not much. Check out `RenderChooseQuestion`, `RenderChooseOption` and `RenderChooseQuery` in [input.go](input.go).

### Extending Output

Output comes with a set of tokens, such as `<success>` or `<error>`, which inject ASCII color codes into the output stream. If you don't want fancy colors, you can just:

``` go
cli := clif.New("my-app", "0.1.0", "Boring output")
cli.SetOutput(clif.NewPlainOutput())
```

To extend or change the fancy style, please modify `clif.DefaultStyles` in [formatter.go](formatter.go).

Patterns & Examples
-------------------

Some patterns I employ which might make sense to others:

### Default options - eg config file

Assuming each/most of your commands require some global config, which needs to have
an optional path. Eg: `my-app do-something --config /path/to/config.yml`.

This can be solved using the `Setup` method of default options:

```go
// Some type for holding config
type MyConfig struct {
    Data map[string]interface{}
}

// init new cli app
cli := clif.New("my-app", "1.2.3", "My app that does something")

// register default option, which fills injection container with config instance
configOpt := clif.NewOption("config", "c", "Path to config file", "/default/config/path.json", true, false).
    SetSetup(function(name, value string) (string, error) {
        if raw, err := ioutil.ReadFile(value); err != nil {
            return "", fmt.Errorf("Could not read config file %s: %s", value, err)
        } else if err = json.Unmarshal(raw, &conf.Data); err != nil {
            return "", fmt.Errorf("Could not unmarshal config file %s: %s", value, err)
        } else if _, ok := conf.Data["name"]; !ok {
            return "", fmt.Errorf("Config %s is missing \"name\"", value)
        } else {
            cli.Register(conf)
            return value, nil
        }
    })

// register command, which uses config instance
cli.New("foo", "Do foo", func(conf *MyConfig) {
    if v, ok := conf.Data["foo"]; ok {
        // ..
    }
}).NewOption("other", "o", "Other option", false, false)
```

### Structure

I like a tidy structure, so I usually have an init-based setup. See [Repo](https://github.com/ukautz/repos) (a tool to keep track of changes in your repos) for an example. Especially the [main.go](https://github.com/ukautz/repos/main/main.go) and the command initialization in [commands.go](https://github.com/ukautz/repos/commands.go).

See also
--------

There are a lot of [other approaches](https://github.com/avelino/awesome-go#command-line) you should have a look at.