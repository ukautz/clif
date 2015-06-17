Go CLI framework
================

Go framework to rapidly develop command line applications.

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
cmd1 := cli.NewCommand("name", "A description", callBackFunction)
cmd2 := cli.NewCommand("other", "Another description", callBackFunction2)
```

The `name` is used from the command line to call the command:

```bash
$ ./app name
$ ./app other
```

### Callback functions

Callback functions can have arbitrary parameters. CLIF uses a small, built-in dependency injection container which allows you to register any kind of object (`struct` or `interface`) beforehand.

CLIF pre-populates the dependency container with a couple of built-in objects:

* The `Output` (formatted output helper, see below)
* The `Input` (input helper, see below)
* The `*Cli` instance itself
* The current `*Command` instance

``` go
package main

import "gopkg.in/ukautz/clif.v0"

type MyFoo struct {}

// ... MyFoo implementation

func callMe(out clif.Output, foo *MyFoo) {
	out.Printf("Foo: %s\n", foo)
}

func main() {
	cli := clif.New("my-app", "My kewl App", "0.8.5").
		New("call", "Call me", callMe).
		Register(new(MyFoo))
	cli.Run()
}
```

If you want to register an `interface`, you need to use the Go'ish way:

``` go
package main

import "gopkg.in/ukautz/clif.v0"

// the interface
type MyFoo interface {
	Bar() string
}

// the struct implementing the interface
type MyBar struct {}

// ... MyBar implementation

func callMe(out clif.Output, foo MyFoo) {
	out.Printf("Foo: %s\n", foo)
}

func main() {
	t := reflect.TypeOf((*MyFoo)(nil)).Elem()
	cli := clif.New("my-app", "My kewl App", "0.8.5").
		New("call", "Call me", callMe).
		RegisterAs(t.String(), new(MyBar))
	cli.Run()
}
```

### Arguments and Options

CLIF can deal with arguments and options. The difference being:

* **Arguments** come after the command name. They are identified by their position.
* **Options** have no fixed position. They are identified by their `--opt-name` (or alias, eg `-O`)

Of course you can use arguments and options at the same time..

#### Arguments

Arguments are additional command line parameters which come after the command name itself.

``` go
cmd := cli.NewCommand("hello", "A description", callBackFunction)
	.NewArgument("name", "Name for greeting", "", true, false)

arg := clif.NewAgument("other", "Something ..", "default", false, true)
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
cmd := cli.NewCommand("hello", "A description", callBackFunction)
	.NewOption("name", "n", "Name for greeting", "", true, false)

arg := clif.NewOption("other", "O", "Something ..", "default", false, true)
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
cmd := cli.NewCommand("hello", "A description", callBackFunction).AddOption(flag)
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

See also
--------

There are a lot of [other approaches](https://github.com/avelino/awesome-go#command-line) you should have a look at.