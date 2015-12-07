// +build ignore

package main

import (
	"fmt"
	"github.com/ukautz/clif"
)

var users = [][]string{
	{
		"Yoda",
		"Very, very old",
		"Like the uber guy",
	},
	{
		"Luke Skywalker",
		"Not that old",
		"A bit, but not that much",
	},
	{
		"Anakin Skywalker",
		"Old dude",
		"He is Lukes father! What do you think?",
	},
	{
		"Super Long Line Super Long Line Super Long Line Super Long Line Super Long Line ",
		"x",
		"x",
	},
	{
		"x",
		"Super Long Line Super Long Line Super Long Line Super Long Line Super Long Line ",
		"x",
	},
	{
		"x",
		"x",
		"Super Long Line Super Long Line Super Long Line Super Long Line Super Long Line ",
	},
	{
		"Super Long Line Super Long Line Super Long Line Super Long Line Super Long Line ",
		"Super Long Line Super Long Line Super Long Line Super Long Line Super Long Line ",
		"x",
	},
	{
		"Super Long Line <info>Super Long<reset> Line Super Long Line Super Long Line Super Long Line ",
		"x",
		"Super Long Line Super Long Line Super Long Line Super Long Line Super Long Line ",
	},
	{
		"x",
		"Super <subline>Long Line Super Long Line Super Long Line Super Long Line Super Long<reset> Line ",
		"Super Long Line Super Long Line Super Long Line Super Long Line Super Long Line ",
	},
	{
		"Super Long Line Super Long Line Super Long Line Super Long Line Super Long Line ",
		"Super Long Line Super Long Line Super Long Line Super Long Line Super Long Line ",
		"Super Long Line Super Long Line Super Long Line Super Long Line Super Long Line ",
	},
}

func printTable(out clif.Output) {
	out.Printf("<headline>Generating table</headline>\n\n")
	headers := []string{"Name", "Age", "Force"}
	table := out.Table(headers)
	for _, user := range users {
		table.AddRow(user)
	}
	fmt.Println(table.Render(0))
}

func main() {
	cli := clif.New("My App", "1.0.0", "An example application").
		New("demo", "Print the progress bar", printTable).
		SetDefaultCommand("demo")
	//cli.SetOutput(clif.NewDebugOutput(os.Stdout))
	cli.Run()
}
