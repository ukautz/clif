package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type commandsSort []*Command

func (this commandsSort) Len() int {
	return len(this)
}

func (this commandsSort) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this commandsSort) Less(i, j int) bool {
	return this[i].Name < this[j].Name
}

type Describer interface {
	Cli(v *Cli) string
	Command(v *Command) string
}

type DefaultDescriber struct {
}

/*
Cli outputs CLI description in long format:

	Name (version)

	Optional description

	Usage:
		prog <command> [<arg> ..] [--opt <val> ..]

	Available Commands:
		foo         This command does foo
		bar         This command does bar
		baz:boing   This command does baz with boing
*/
func (this *DefaultDescriber) Cli(v *Cli) string {

	// name + version
	line := v.Name
	if v.Version != "" {
		line += " (" + v.Version + ")"
	}
	lines := []string{line}

	// description
	if v.Description != "" {
		lines = append(lines, v.Description+"\n")
	}

	// usage
	prog := filepath.Base(os.Args[0])
	lines = append(lines, fmt.Sprintf("Usage:\n\t%s <command> [<arg> ..] [--opt <val> ..]\n", prog))

	// commands
	lines = append(lines, "Available commands:")
	max := 0
	commands := make([]*Command, len(v.Commands))
	i := 0
	for _, cmd := range v.Commands {
		commands[i] = cmd
		i++
		if l := len(cmd.Name); l > max {
			max = l
		}
	}
	sort.Sort(commandsSort(commands))
	for _, cmd := range commands {
		lines = append(lines, fmt.Sprintf("\t%-"+fmt.Sprintf("%d", max)+"s  %s", cmd.Name, cmd.Usage))
	}

	return strings.Join(lines, "\n")
}

/*
Returns:

	(Description)

	Usage:
		command <arg> [<arg> ..] --opt <val> (--opt <val> ..)

	Arguments:
		foo  Foo is needed (required)
		bar  Bar is optional (multiple)

	Options:
		--bla (-f)  Something?
		--blub      Morething? (required) (multiple) (default=..)
*/
func (this *DefaultDescriber) Command(v *Command) string {
	lines := []string{}
	if v.Description != "" {
		lines = append(lines, []string{v.Description, ""}...)
	} else if v.Usage != "" {
		lines = append(lines, []string{v.Usage, ""}...)
	}

	lines = append(lines, "Usage:")
	usage := []string{v.Name}
	args := make([][]string, 0)
	argMax := 0
	opts := make([][]string, 0)
	optMax := 0
	for _, p := range v.Arguments {
		var short string
		usg := p.Usage
		if p.Required {
			short = fmt.Sprintf("<%s>", p.Name)
			usg += " (req)"
		} else {
			short = fmt.Sprintf("[%s]", p.Name)
		}
		if p.Multiple {
			short = "(" + short + " ...)"
			usg += " (mult)"
		}
		if p.Default != "" {
			usg += fmt.Sprintf(" (default: \"%s\")", p.Default)
		}
		if l := len(p.Name); l > argMax {
			argMax = l
		}
		usage = append(usage, short)
		args = append(args, []string{p.Name, usg})
	}
	for _, p := range v.Options {
		short := fmt.Sprintf("--%s", p.Name)
		if p.Alias != "" {
			short += "|-" + p.Alias
		}
		if !p.Flag {
			short += " <val>"
		}
		long := short
		usg := p.Usage
		if !p.Required {
			short = "(" + short + ")"
		} else {
			usg += " (req)"
		}
		if p.Multiple {
			short = "(" + short + " ...)"
			usg += " (mult)"
		}
		if p.Default != "" {
			usg += fmt.Sprintf(" (default: \"%s\")", p.Default)
		}
		if l := len(long); l > optMax {
			optMax = l
		}
		usage = append(usage, short)
		opts = append(opts, []string{long, usg})
	}
	lines = append(lines, "\t"+strings.Join(usage, " "))
	lines = append(lines, "")

	if len(args) > 0 {
		lines = append(lines, "Arguments:")
		for _, l := range args {
			lines = append(lines, fmt.Sprintf("\t%-"+fmt.Sprintf("%d", argMax)+"s  %s", l[0], l[1]))
		}
		lines = append(lines, "")
	}

	if len(opts) > 0 {
		lines = append(lines, "Options:")
		for _, l := range opts {
			lines = append(lines, fmt.Sprintf("\t%-"+fmt.Sprintf("%d", optMax)+"s  %s", l[0], l[1]))
		}
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}
