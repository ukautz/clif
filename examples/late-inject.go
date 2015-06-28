// +build ignore

// Example on how to late-inject objects into container registry, based on default option
package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/ukautz/clif.v0"
	"io/ioutil"
)

var description = `
This application demonstrates how to use default parameters and utilize late-injection
to solve the common case when each command requires a config file.
`

type MyConfig struct {
	Data map[string]string
}

func main() {

	// init cli app
	cli := clif.New("my-app", "Demo for using config", "1.2.3").SetDescription(description)

	// add a default option, which registers a new object in the injection container
	cli.AddDefaultOptions(
		clif.NewOption("config", "c", "Path to config", "fixtures/config.json", true, false).SetSetup(func(name, value string) (string, error) {
			conf := &MyConfig{make(map[string]string)}
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
		}),
	)

	// add command which uses the late injected configf
	cli.Add(clif.NewCommand("xxx", "Call xxx", func(c *clif.Command, foo *MyConfig, out clif.Output) {
		out.Printf("Hello there: <success>%s<reset>\n", foo.Data["name"])
	}))

	// add another ocmmand, using the config as well
	cli.Add(clif.NewCommand("yyy", "Call yyy", func(c *clif.Command, foo *MyConfig, out clif.Output) {
		out.Printf("Hello there: <success>%s<reset>\n", foo.Data["name"])
	}))

	cli.Run()
}
