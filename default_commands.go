package clif

import "fmt"

// NewHelpCommand returns the default help command
func NewHelpCommand() *Command {
	return NewCommand("help", "Show this help", func(o *Command, out Output) error {
		if n := o.Argument("command").String(); n != "" {
			if cmd, ok := o.Cli.Commands[n]; ok {
				out.Printf(DescribeCommand(cmd))
			} else {
				out.Printf(DescribeCli(o.Cli))
				return fmt.Errorf("Unknown command \"%s\"", n)
			}
		} else {
			out.Printf(DescribeCommand(o))
		}
		return nil
	}).NewArgument("command", "Command to show help for", "", false, false)
}

// NewListCommand returns the default help command
func NewListCommand() *Command {
	return NewCommand("list", "List all available commands", func(c *Cli, Command, out Output) {
		out.Printf(DescribeCli(c))
	})
}
