// +build ignore

package main

import (
	"fmt"
	"github.com/ukautz/clif"
)

func printTable1(c *clif.Command, out clif.Output) {
	headers := []string{"H1", "H2", "H3"}
	var table *clif.Table
	if c.Option("open").Bool() {
		table = out.Table(headers, clif.OpenTableStyle)
	} else {
		table = out.Table(headers)
	}
	table.AddRows([][]string{
		[]string{"foo", "bar", "baz"},
		[]string{"yadda", "yadda", "yadda"},
		[]string{"Some crazy multi line content + Some crazy multi line content + Some crazy multi line content", "yadda", "yadda"},
		[]string{"yadda", "Some crazy multi line content + Some crazy multi line content + Some crazy multi line content", "yadda"},
		[]string{"yadda", "yadda", "Some crazy multi line content + Some crazy multi line content + Some crazy multi line content"},
		[]string{"Some <info>crazy multi line content + Some crazy multi line content + Some crazy<reset> multi line content", "yadda", "yadda"},
		[]string{"yadda", "Some <info>crazy multi line content + Some crazy multi line content + Some crazy<reset> multi line content", "yadda"},
		[]string{"yadda", "yadda", "Some <info>crazy multi line content + Some crazy multi line content + Some crazy<reset> multi line content"},
	})
	fmt.Println(table.Render(c.Option("render-width").Int()))
}

func printTable2(c *clif.Command, out clif.Output) {
	out.Printf("<headline>Generating table</headline>\n\n")
	headers := []string{"Name", "Age", "Force"}
	var table *clif.Table
	if c.Option("open").Bool() {
		table = out.Table(headers, clif.OpenTableStyleLight)
	} else {
		table = out.Table(headers)
	}
	users := [][]string{
		{
			"<important>Yoda<reset>",
			"Very, very old",
			"Like the uber guy",
		},
		{
			"<important>Luke Skywalker<reset>",
			"Not that old",
			"A bit, but not that much",
		},
		{
			"<important>Anakin Skywalker<reset>",
			"Old dude",
			"He is Lukes father! Was kind of stronger in 1-3, but still failed to" +
				" kill Jar Jar Binks. Not even tried, though. What's with that?",
		},
	}
	for _, user := range users {
		table.AddRow(user)
	}
	fmt.Println(table.Render(c.Option("render-width").Int()))
}

func main() {
	cli := clif.New("My App", "1.0.0", "An example application").
		New("demo1", "Print demo table 1", printTable1).
		New("demo2", "Print demo table 2", printTable2).
		AddDefaultOptions(clif.NewFlag("open", "O", "Use open table style", false)).
		AddDefaultOptions(clif.NewOption("render-width", "w", "Render width of the table", fmt.Sprintf("%d", clif.TermWidthCurrent), false, false)).
		SetDefaultCommand("demo1")
	//cli.SetOutput(clif.NewDebugOutput(os.Stdout))
	cli.Run()
}
