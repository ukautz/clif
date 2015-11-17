// +build ignore

package main

import (
	"github.com/ukautz/clif"
	"os"
	"fmt"
)

var forceUsers = [][]string{
	{
		"Yoda",
		"Very, very old",
		"Like the uber guy",
	},
	{
		"Luke Skywalker",
		"<info>Not that old</reset>",
		"A bit, but not that much",
	},
	{
		"Anakin Skywalker",
		"Old dude",
		"He is Lukes father! What do you think?",
	},
}

func printTable(out clif.Output) {
	out.Printf("<headline>Generating table</headline>\n\n")
	headers := []string{"Name", "Age", "Force"}
	table := out.Table(headers)
	table.Style.ContentRenderer = func(str string) string {
		fmt.Printf("\nOUT: \"%s\"\n", str)
		return out.Sprintf(str)
	}
	for _, row := range forceUsers {
		table.AddRow(row)
	}
	fmt.Println(table.Render(80))
}

func main() {
	cli := clif.New("My App", "1.0.0", "An example application").
		New("demo", "Print the progress bar", printTable).
		SetDefaultCommand("demo")
	cli.SetOutput(clif.NewDebugOutput(os.Stdout))
	cli.Run()
}
