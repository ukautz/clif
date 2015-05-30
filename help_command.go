package cli

func NewHelpCommand() *Command {
	return NewCommand("help", "Show this help", func(c *Cli, o *Command) {
		if n := o.Argument("command").String(); n != "" {
			if cmd, ok := o.Cli.Commands[n]; ok {
				c.Output.Printf(o.Cli.Describer.Command(cmd))
			} else {
				c.Output.Printf(o.Cli.Describer.Cli(o.Cli))
				c.Output.Printf("\n\nUnknown command \"%s\"\n", n)
			}
		} else {
			c.Output.Printf(o.Cli.Describer.Command(o))
		}
	}).NewArgument("command", "Command to show help for", "", false, false)
}
